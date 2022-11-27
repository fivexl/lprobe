package main

import (
	"lprobe/grpc-health-probe"
	"strconv"
)

//nolint:all
func grpc_healthcheck(port int) (string) {
	status, code := grpchealthprobe.Grpchealthprobe("127.0.0.1:" + strconv.Itoa(port))
	if code != 0 {
		// Error status
		return status
	} else {
		return ""
	}
}
