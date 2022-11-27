package main

import (
	"fmt"
	"os"
	"flag"
)


func main() {

	port := flag.Int("port", 8080, "The port number")
	endpoint := flag.String("endpoint", "/", "Endpoint for probe")
	mode := flag.String("mode", "http", "HTTP or gRPC mode for checks")
	flag.Parse()
	modeValue := *mode
	portValue := *port
	if modeValue == "http" {
		// HTTP check
		endpointValue := *endpoint
		_, err := http_healthcheck(portValue, endpointValue)
		if err != nil {
			fmt.Printf("Error: %v", err)
			os.Exit(1)
		}
	} else if modeValue == "grpc" {
		// gRPC check
		status := grpc_healthcheck(portValue)
		if status != "" {
			fmt.Printf("Error: %v", status)
			os.Exit(1)
		}
		os.Exit(0)
	} else {
		// unkown check
		fmt.Printf("Error: mode is not supported")
		os.Exit(1)
	}
}
