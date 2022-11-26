package main

import (
	"fmt"
	"net/http"
	"strconv"
)

func http_healthcheck(port int, endpoint string) (*http.Response, error) {
	//TODO if we get endpoint without / , add it 
	return http.Get(fmt.Sprintf("http://127.0.0.1:%s%s", strconv.Itoa(port), endpoint))
}
