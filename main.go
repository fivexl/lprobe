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
	https := flag.Bool("https", false, "HTTPS or Not. Used for HTTP check.")
	flag.Parse()
	modeValue := *mode
	portValue := *port
	if modeValue == "http" {
		// HTTP check
		endpointValue := *endpoint
		httpsValue := *https
		var protocol string
		if httpsValue {
			protocol = "https"
		} else {
			protocol = "http"
		}
		_, err := http_healthcheck(protocol, portValue, endpointValue)
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
