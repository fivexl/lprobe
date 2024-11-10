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
	"strconv"
	"context"
	"log"
	"time"

	"github.com/spiffe/go-spiffe/v2/spiffetls/tlsconfig"
	"github.com/spiffe/go-spiffe/v2/workloadapi"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/alts"
	"google.golang.org/grpc/credentials/insecure"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

//nolint:all
func grpcHealthCheck() (string) {
	status, code := grpchealthprobe(getAddr() + ":" + strconv.Itoa(flPort))
	if code != 0 {
		// Error status
		return status
	}
	return ""
}

func grpchealthprobe(flAddr string) (string, int) {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()


	opts := []grpc.DialOption{
		grpc.WithUserAgent(flUserAgent),
		grpc.WithTimeout(flConnTimeout),
	}
	if flTLS && flSPIFFE {
		log.Printf("-tls and -spiffe are mutually incompatible")
		return "ERR", StatusInvalidArguments
	}
	if flTLS {
		creds, _ , err := buildCredentials(flTLSNoVerify, flTLSCACert, flTLSClientCert, flTLSClientKey, flTLSServerName)
		if err != nil {
			log.Printf("failed to initialize tls credentials. error=%v", err)
			return "ERR", StatusInvalidArguments
		}
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else if flALTS {
		creds := alts.NewServerCreds(alts.DefaultServerOptions())
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else if flSPIFFE {
		spiffeCtx, cancel := context.WithTimeout(ctx, flRPCTimeout)
		defer cancel()
		source, err := workloadapi.NewX509Source(spiffeCtx)
		if err != nil {
			log.Printf("failed to initialize tls credentials with spiffe. error=%v", err)
			return "ERR", StatusSpiffeFailed
		}
		if flVerbose {
			svid, err := source.GetX509SVID()
			if err != nil {
				log.Fatalf("error getting x509 svid: %+v", err)
			}
			log.Printf("SPIFFE Verifiable Identity Document (SVID): %q", svid.ID)
		}
		creds := credentials.NewTLS(tlsconfig.MTLSClientConfig(source, source, tlsconfig.AuthorizeAny()))
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	if flGZIP {
		opts = append(opts,
			grpc.WithCompressor(grpc.NewGZIPCompressor()), 		//nolint:all 
			grpc.WithDecompressor(grpc.NewGZIPDecompressor()),	//nolint:all
		)
	}

	if flVerbose {
		log.Print("establishing connection")
	}
	connStart := time.Now()
	conn, err := grpc.NewClient(
		flAddr,
		opts...,
	)
	if err != nil {
		if err == context.DeadlineExceeded {
			log.Printf("timeout: failed to connect service %q within %v", flAddr, flConnTimeout)
		} else {
			log.Printf("error: failed to connect service at %q: %+v", flAddr, err)
		}
		return "ERR", StatusConnectionFailure
	}
	connDuration := time.Since(connStart)
	defer conn.Close()
	if flVerbose {
		log.Printf("connection established (took %v)", connDuration)
	}

	rpcStart := time.Now()
	rpcCtx, rpcCancel := context.WithTimeout(ctx, flRPCTimeout)
	defer rpcCancel()
	rpcCtx = metadata.NewOutgoingContext(rpcCtx, flRPCHeaders.MD)
	resp, err := healthpb.NewHealthClient(conn).Check(rpcCtx,
		&healthpb.HealthCheckRequest{
			Service: flService})
	if err != nil {
		if stat, ok := status.FromError(err); ok && stat.Code() == codes.Unimplemented {
			log.Printf("error: this server does not implement the grpc health protocol (grpc.health.v1.Health): %s", stat.Message())
		} else if stat, ok := status.FromError(err); ok && stat.Code() == codes.DeadlineExceeded {
			log.Printf("timeout: health rpc did not complete within %v", flRPCTimeout)
		} else {
			log.Printf("error: health rpc failed: %+v", err)
		}
		return "ERR", StatusRPCFailure
	}
	rpcDuration := time.Since(rpcStart)

	if resp.GetStatus() != healthpb.HealthCheckResponse_SERVING {
		return resp.GetStatus().String(), StatusUnhealthy
	}
	if flVerbose {
		log.Printf("time elapsed: connect=%v rpc=%v", connDuration, rpcDuration)
	}
	return  resp.GetStatus().String(), 0
}
