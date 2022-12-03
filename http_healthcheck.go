package main

import (
	"fmt"
	"net/http"
	"strconv"
)

func httpHealthCheck(port int, endpoint string) (*http.Response, error) {
	//TODO if we get endpoint without / , add it 
	//TODO add TLS support
	var protocol string
	if flTLS {
		protocol = "https"
	} else {
		protocol = "http"
	}
	return http.Get(fmt.Sprintf("%s://%s:%s%s", protocol, LocalAddress, strconv.Itoa(port), endpoint))
}
