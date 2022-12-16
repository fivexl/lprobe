package main

import (
	"net/http"
	"strconv"
	"strings"
	"crypto/tls"
)

func httpHealthCheck() (*http.Response, error) {

	var	endpoint string
	var	protocol string

	if strings.HasPrefix(flEndpoint, "/") {
		endpoint = flEndpoint
	} else {
		endpoint = "/" + flEndpoint
	}

	if flTLS {
		protocol = "https"
	} else {
		protocol = "http"
	}

    httpTransport := &http.Transport{
        TLSClientConfig: &tls.Config{InsecureSkipVerify: flTLSNoVerify},
    }

	httpClient := &http.Client{
		Transport: httpTransport,
		Timeout: flConnTimeout,
	}

	url := protocol + "://" + getAddr() + ":" + strconv.Itoa(flPort) + endpoint
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", flUserAgent)
	res, err := httpClient.Do(req)
	return res, err
}
