package main

import (
	"fmt"
	"os"
)

func main() {

	switch flMode {
	case "http":
		// HTTP check
		err := httpHealthCheck()
		if err != nil {
			fmt.Printf("Error: %v", err)
			os.Exit(1)
		}
	case "grpc":
		// gRPC check
		status := grpcHealthCheck()
		if status != "" {
			fmt.Printf("Error: %v", status)
			os.Exit(1)
		}
		os.Exit(0)
	default:
		// unkown check
		fmt.Printf("Error: Unsupported -mode. Please use one of %v", getSupportedModes())
		os.Exit(1)
	}
}
