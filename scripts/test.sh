#!/usr/bin/env bash

# set -e
rm -rf ./lprobe
go build

## Prepull latest images
docker pull nginx
docker pull grpc/java-example-hostname

## Run docker containers
docker run --rm -p 8080:80 -d --name nginx-lprobe-test nginx 
docker run --rm -p 8081:50051 -d --name grpc-lprobe-test grpc/java-example-hostname
echo "Wait 5s"
sleep 5

### HTTP Check Test
./lprobe -mode=http -port=8080 -endpoint=/ -http-codes=200-298,299 -v
if [ "$?" != 0 ]; then
    echo "HTTP test failed"
    docker stop nginx-lprobe-test
    exit 1
fi

### HTTP IPv6 Check Test
./lprobe -mode=http -port=8080 -endpoint=/ -ipv6 -v
if [ "$?" != 0 ]; then
    echo "HTTP IPv6 test failed"
    docker stop nginx-lprobe-test
    exit 1
fi

### gRPC Check Test
./lprobe -mode=grpc -port=8081 -v
if [ "$?" != 0 ]; then
    echo "gRPC test failed"
    docker stop grpc-lprobe-test
    exit 1
fi

### gRPC IPv6 Check Test
./lprobe -mode=grpc -port=8081 -ipv6 -v
if [ "$?" != 0 ]; then
    echo "gRPC IPv6 test failed"
    docker stop grpc-lprobe-test
    exit 1
fi


### FAIL HTTP Check Test
./lprobe -mode=http -port=7777 -endpoint=/ -v
if [ "$?" != 1 ]; then
    echo "FAIL HTTP test failed"
    docker stop nginx-lprobe-test
    exit 1
fi

### FAIL gRPC Check Test
./lprobe -mode=grpc -port=7777 -v
if [ "$?" != 1 ]; then
    echo "FAIL gRPC test failed"
    docker stop grpc-lprobe-test
    exit 1
fi

### URL HTTP Check Test
./lprobe -url http://127.0.0.1:8080/ -v
if [ "$?" != 0 ]; then
    echo "URL HTTP test failed"
    docker stop nginx-lprobe-test
    exit 1
fi

### URL HTTP with custom path Check Test
./lprobe -url http://127.0.0.1:8080/ -v
if [ "$?" != 0 ]; then
    echo "URL HTTP path test failed"
    docker stop nginx-lprobe-test
    exit 1
fi

### FAIL URL HTTP Check Test
./lprobe -url http://127.0.0.1:7777/ -v
if [ "$?" != 1 ]; then
    echo "FAIL URL HTTP test failed"
    docker stop nginx-lprobe-test
    exit 1
fi

## Stop docker containers
docker stop nginx-lprobe-test
docker stop grpc-lprobe-test

echo "All good"
exit 0