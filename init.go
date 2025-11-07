// Copyright 2018 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/url"
	"os"
	"strings"
	"time"
	"unicode"

	"golang.org/x/exp/slices"
	"google.golang.org/grpc/metadata"
)

var (
	flMode			string
	flIPv6			bool
	flPort 			int
	flUserAgent     string	
	flEndpoint 		string
	flHTTPCodes     string  
	flService       string
	flConnTimeout   time.Duration
	flRPCHeaders    = rpcHeaders{MD: make(metadata.MD)}
	flRPCTimeout    time.Duration
	flTLS           bool
	flTLSNoVerify   bool
	flTLSCACert     string
	flTLSClientCert string
	flTLSClientKey  string
	flTLSServerName string
	flALTS          bool
	flVerbose       bool
	flGZIP          bool
	flSPIFFE        bool
	flURL           string
	flHost          string
)

const (
	// LocalAddress to call 
	LocalAddress = "127.0.0.1"
	// LocalAddress6 IPv6 to call 
	LocalAddress6 = "[::1]"
	// AWS metadata IPv4 address
	AWSMetadataIPv4 = "169.254.169.254"
	// AWS EKS metadata IPv4 address
	AWSEKSIPv4 = "169.254.170.2"
	// AWS metadata IPv6 address
	AWSMetadataIPv6 = "fd00:ec2::254"
	// StatusInvalidArguments indicates specified invalid arguments.
	StatusInvalidArguments = 1
	// StatusConnectionFailure indicates connection failed.
	StatusConnectionFailure = 2
	// StatusRPCFailure indicates rpc failed.
	StatusRPCFailure = 3
	// StatusUnhealthy indicates rpc succeeded but indicates unhealthy service.
	StatusUnhealthy = 4
	// StatusSpiffeFailed indicates failure to retrieve credentials using spiffe workload API
	StatusSpiffeFailed = 20
)

func getSupportedModes() []string {
	return []string{"http", "grpc"}
}

func validateHostname(hostname string) error {
	if hostname == "" {
		return nil
	}

	// Check if it's already an IP address
	if ip := net.ParseIP(hostname); ip != nil {
		return validateIPAddress(ip)
	}

	// Resolve hostname to IP addresses
	ips, err := net.LookupIP(hostname)
	if err != nil {
		return fmt.Errorf("failed to resolve hostname %q: %v", hostname, err)
	}

	// Check if any resolved IP is allowed
	for _, ip := range ips {
		if err := validateIPAddress(ip); err != nil {
			return fmt.Errorf("hostname %q resolves to disallowed IP %s: %v", hostname, ip, err)
		}
	}

	return nil
}

func validateIPAddress(ip net.IP) error {
	// Check IPv4 addresses
	if ip.To4() != nil {
		if ip.String() == LocalAddress || ip.String() == AWSMetadataIPv4 || ip.String() == AWSEKSIPv4 {
			return nil
		}
		return fmt.Errorf("IPv4 address %s is not allowed (only localhost and AWS metadata endpoints permitted)", ip)
	}

	// Check IPv6 addresses
	if ip.To16() != nil {
		// Remove zone index if present (e.g., ::1%eth0)
		ipWithoutZone := net.IP(ip.To16())
		if ipWithoutZone.String() == "::1" || ipWithoutZone.String() == AWSMetadataIPv6 {
			return nil
		}
		return fmt.Errorf("IPv6 address %s is not allowed (only localhost and AWS metadata endpoints permitted)", ip)
	}

	return fmt.Errorf("invalid IP address format: %s", ip)
}

func getAddr() string {
	// If URL is provided, use hostname from URL (already validated)
	if flURL != "" && flHost != "" {
		return flHost
	}
	// Otherwise use localhost
	if flIPv6 {
		return LocalAddress6
	}
	return LocalAddress
}

func init() {
	flagSet := flag.NewFlagSet("", flag.ContinueOnError)
	log.SetFlags(0)
	// core settings
	flagSet.StringVar(&flMode, "mode", "http", "Select mode: http, grpc (default: http)")
	flagSet.BoolVar(&flIPv6, "ipv6", false, "Use IPv6 ::1 address (default: false)")
	flagSet.IntVar(&flPort, "port", 8080, "port number to check (defaut 8080)")
	flagSet.StringVar(&flUserAgent, "user-agent", "lprobe", "user-agent header value of health check requests")
	// HTTP settings
	flagSet.StringVar(&flEndpoint, "endpoint", "/", "HTTP endpoint (default: /)")
	flagSet.StringVar(&flHTTPCodes, "http-codes", "200-299", "comma-separated list of expected HTTP status codes (default: 200-299) with optional ranges (e.g. 200,201,202-204,301-303")
	// gRPC settings
	flagSet.StringVar(&flService, "service", "", "service name to check (default: \"\")")
	// timeouts
	flagSet.DurationVar(&flConnTimeout, "connect-timeout", time.Second, "timeout for establishing connection")
	// headers
	flagSet.Var(&flRPCHeaders, "rpc-header", "additional RPC headers in 'name: value' format. May specify more than one via multiple flags.")
	flagSet.DurationVar(&flRPCTimeout, "rpc-timeout", time.Second, "timeout for health check rpc")
	// tls settings
	flagSet.BoolVar(&flTLS, "tls", false, "use TLS (default: false, INSECURE plaintext transport)")
	flagSet.BoolVar(&flTLSNoVerify, "tls-no-verify", false, "(with -tls) don't verify the certificate (INSECURE) presented by the server (default: false)")
	flagSet.StringVar(&flTLSCACert, "tls-ca-cert", "", "(with -tls, optional) file containing trusted certificates for verifying server")
	flagSet.StringVar(&flTLSClientCert, "tls-client-cert", "", "(with -tls, optional) client certificate for authenticating to the server (requires -tls-client-key)")
	flagSet.StringVar(&flTLSClientKey, "tls-client-key", "", "(with -tls) client private key for authenticating to the server (requires -tls-client-cert)")
	flagSet.StringVar(&flTLSServerName, "tls-server-name", "", "(with -tls) override the hostname used to verify the server certificate")
	flagSet.BoolVar(&flALTS, "alts", false, "use ALTS (default: false, INSECURE plaintext transport)")
	flagSet.BoolVar(&flVerbose, "v", false, "verbose logs")
	flagSet.BoolVar(&flGZIP, "gzip", false, "use GZIPCompressor for requests and GZIPDecompressor for response (default: false)")
	flagSet.BoolVar(&flSPIFFE, "spiffe", false, "use SPIFFE to obtain mTLS credentials")
	flagSet.StringVar(&flURL, "url", "", "Complete URL to check (e.g., http://example.com:8080/path)")

	err := flagSet.Parse(os.Args[1:])
	if err != nil {
		os.Exit(StatusInvalidArguments)
	}

	argError := func(s string, v ...interface{}) {
		log.Printf("error: "+s, v...)
		os.Exit(StatusInvalidArguments)
	}

	// Parse URL if provided
	if flURL != "" {
		parsedURL, err := url.Parse(flURL)
		if err != nil {
			argError("invalid URL format: %v", err)
		}

		// Set mode based on scheme
		if parsedURL.Scheme == "https" {
			flTLS = true
		} else if parsedURL.Scheme != "http" {
			argError("unsupported URL scheme: %s (only http and https are supported)", parsedURL.Scheme)
		}

		// Override port if specified in URL
		if parsedURL.Port() != "" {
			fmt.Sscanf(parsedURL.Port(), "%d", &flPort)
		} else if parsedURL.Scheme == "https" {
			flPort = 443
		} else if parsedURL.Scheme == "http" {
			flPort = 80
		}

		// Override endpoint if path is specified
		if parsedURL.Path != "" {
			flEndpoint = parsedURL.Path
		}

		// Set hostname from URL
		if parsedURL.Hostname() != "" {
			hostname := parsedURL.Hostname()
			if err := validateHostname(hostname); err != nil {
				argError("hostname validation failed: %v", err)
			}
			flHost = hostname
		}

		// Set server name for TLS verification if hostname is available
		if flTLS && parsedURL.Hostname() != "" {
			flTLSServerName = parsedURL.Hostname()
		}
	}

	if !slices.Contains(getSupportedModes(), flMode)  {
		argError("Unsupported -mode. Please use one of %v", getSupportedModes())
	}
	if flConnTimeout <= 0 {
		argError("-connect-timeout must be greater than zero (specified: %v)", flConnTimeout)
	}
	if flRPCTimeout <= 0 {
		argError("-rpc-timeout must be greater than zero (specified: %v)", flRPCTimeout)
	}
	if flALTS && flSPIFFE {
		argError("-alts and -spiffe are mutually incompatible")
	}
	if flTLS && flALTS {
		argError("cannot specify -tls with -alts")
	}
	if !flTLS && flTLSNoVerify {
		argError("specified -tls-no-verify without specifying -tls")
	}
	if !flTLS && flTLSCACert != "" {
		argError("specified -tls-ca-cert without specifying -tls")
	}
	if !flTLS && flTLSClientCert != "" {
		argError("specified -tls-client-cert without specifying -tls")
	}
	if !flTLS && flTLSServerName != "" {
		argError("specified -tls-server-name without specifying -tls")
	}
	if flTLSClientCert != "" && flTLSClientKey == "" {
		argError("specified -tls-client-cert without specifying -tls-client-key")
	}
	if flTLSClientCert == "" && flTLSClientKey != "" {
		argError("specified -tls-client-key without specifying -tls-client-cert")
	}
	if flTLSNoVerify && flTLSCACert != "" {
		argError("cannot specify -tls-ca-cert with -tls-no-verify (CA cert would not be used)")
	}
	if flTLSNoVerify && flTLSServerName != "" {
		argError("cannot specify -tls-server-name with -tls-no-verify (server name would not be used)")
	}

	if flVerbose {
		log.Printf("parsed options:")
		log.Printf("> conn_timeout=%v rpc_timeout=%v", flConnTimeout, flRPCTimeout)
		if flRPCHeaders.Len() > 0 {
			log.Printf("> headers: %s", flRPCHeaders)
		}
		log.Printf("> tls=%v", flTLS)
		if flTLS {
			log.Printf("  > no-verify=%v ", flTLSNoVerify)
			log.Printf("  > ca-cert=%s", flTLSCACert)
			log.Printf("  > client-cert=%s", flTLSClientCert)
			log.Printf("  > client-key=%s", flTLSClientKey)
			log.Printf("  > server-name=%s", flTLSServerName)
		}
		log.Printf("> alts=%v", flALTS)
		log.Printf("> spiffe=%v", flSPIFFE)
	}
}

type rpcHeaders struct{ metadata.MD }

func (s *rpcHeaders) String() string { return fmt.Sprintf("%v", s.MD) }

func (s *rpcHeaders) Set(value string) error {
	parts := strings.SplitN(value, ":", 2)
	if len(parts) != 2 {
		return fmt.Errorf("invalid RPC header, expected 'key: value', got %q", value)
	}
	trimmed := strings.TrimLeftFunc(parts[1], unicode.IsSpace)
	s.Append(parts[0], trimmed)
	return nil
}