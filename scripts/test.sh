#!/usr/bin/env bash

set -e
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
./lprobe -mode=http -port=8080 -endpoint=/
if [ "$?" != 0 ]; then
    echo "HTTP test failed"
    docker stop nginx-lprobe-test
    exit 1
fi

### HTTP IPv6 Check Test
./lprobe -mode=http -port=8080 -endpoint=/ -ipv6
if [ "$?" != 0 ]; then
    echo "HTTP IPv6 test failed"
    docker stop nginx-lprobe-test
    exit 1
fi

### gRPC Check Test
./lprobe -mode=grpc -port=8081 
if [ "$?" != 0 ]; then
    echo "gRPC test failed"
    docker stop grpc-lprobe-test
    exit 1
fi

### gRPC IPv6 Check Test
./lprobe -mode=grpc -port=8081 -ipv6
if [ "$?" != 0 ]; then
    echo "gRPC IPv6 test failed"
    docker stop grpc-lprobe-test
    exit 1
fi

## Stop docker containers
docker stop nginx-lprobe-test
docker stop grpc-lprobe-test

echo "All good"
exit 0