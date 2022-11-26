package main

import (
	"fmt"
	"os"
	"flag"
)


func main() {

	port := flag.Int("port", 8080, "The port number")
	endpoint := flag.String("endpoint", "/", "Endpoint for probe")
	flag.Parse()
	portValue := *port
	endpointValue := *endpoint

	_, err := http_healthcheck(portValue, endpointValue)
	if err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}
