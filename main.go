package main

import (
	"fmt"
	"os"
)

func main() {

	if flMode == "http" {
		// HTTP check
		_, err := httpHealthCheck(flPort, flEndpoint)
		if err != nil {
			fmt.Printf("Error: %v", err)
			os.Exit(1)
		}
	} else if flMode == "grpc" {
		// gRPC check
		status := grpcHealthCheck(flPort)
		if status != "" {
			fmt.Printf("Error: %v", status)
			os.Exit(1)
		}
		os.Exit(0)
	} else {
		// unkown check
		fmt.Printf("Error: Unsupported -mode. Please use one of %v", getSupportedModes())
		os.Exit(1)
	}
}
