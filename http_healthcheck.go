package main

import (
	"fmt"
	"net/http"
	"strconv"
)

func http_healthcheck(protocol string, port int, endpoint string) (*http.Response, error) {
	//TODO if we get endpoint without / , add it 
	//TODO add TLS support
	return http.Get(fmt.Sprintf("%s://127.0.0.1:%s%s", protocol, strconv.Itoa(port), endpoint))
}
