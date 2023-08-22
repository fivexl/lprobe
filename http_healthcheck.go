package main

import (
	"net/http"
	"strconv"
	"log"	
	"strings"
	"fmt"
	"sort"
)

func httpHealthCheck() (error) {

	var	endpoint string
	var	protocol string
	var validHTTPCodes map[int]bool
	var httpTransport *http.Transport

	if strings.HasPrefix(flEndpoint, "/") {
		endpoint = flEndpoint
	} else {
		endpoint = "/" + flEndpoint
	}

	validHTTPCodes = make(map[int]bool)
	ranges := strings.Split(flHTTPCodes, ",")
	for _, r := range ranges {
		if strings.Contains(r, "-") {
			parts := strings.Split(r, "-")
			start, _ := strconv.Atoi(parts[0])
			end, _ := strconv.Atoi(parts[1])
			for i := start; i <= end; i++ {
				validHTTPCodes[i] = true
			}
		} else {
			code, _ := strconv.Atoi(r)
			validHTTPCodes[code] = true
		}
	}

	// print validHTTPCodes map for debug
	if flVerbose {
		validCodes := make([]int, 0, len(validHTTPCodes))
		for k := range validHTTPCodes {
			validCodes = append(validCodes, k)
		}
		sort.Ints(validCodes) // Sorting the valid codes for readability
		log.Printf("Valid HTTP Codes: %v\n", validCodes)
	}

	if flTLS {
		protocol = "https"
		_, creds, err := buildCredentials(flTLSNoVerify, flTLSCACert, flTLSClientCert, flTLSClientKey, flTLSServerName)
		if err != nil {
			log.Printf("failed to initialize tls credentials. error=%v", err)
			return err
		}
		httpTransport = &http.Transport{
			TLSClientConfig: creds,
		}
	} else {
		protocol = "http"
		httpTransport = &http.Transport{}
	}

	httpClient := &http.Client{
		Transport: httpTransport,
		Timeout: flConnTimeout,
	}

	// initialize http request
	url := protocol + "://" + getAddr() + ":" + strconv.Itoa(flPort) + endpoint
	req, reqErr := http.NewRequest("GET", url, nil)
	if reqErr != nil {
		log.Printf("failed to initialize http request. error=%v", reqErr)
		return reqErr
	}
	req.Header.Set("User-Agent", flUserAgent)

	// execute http request
	httpResponse, respErr := httpClient.Do(req)
	if respErr != nil {
		log.Printf("failed to execute http request. error=%v", respErr)
		return respErr
	}
	defer httpResponse.Body.Close()

	// check http response code
	if !validHTTPCodes[httpResponse.StatusCode] {
		log.Printf("HTTP request returned status %v", httpResponse.StatusCode)
		return fmt.Errorf("unexpected HTTP status code: %v", httpResponse.StatusCode)
	}

	return nil

}
