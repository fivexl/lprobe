package main

import (
	"net/http"
	"strconv"
	"log"	
	"strings"
)

func httpHealthCheck() (error) {

	var	endpoint string
	var	protocol string
	var httpTransport *http.Transport

	if strings.HasPrefix(flEndpoint, "/") {
		endpoint = flEndpoint
	} else {
		endpoint = "/" + flEndpoint
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

	url := protocol + "://" + getAddr() + ":" + strconv.Itoa(flPort) + endpoint
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", flUserAgent)
	_, err := httpClient.Do(req)
	return err
}
