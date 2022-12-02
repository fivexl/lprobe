#!/usr/bin/env bash

set -e
rm -rf ./lprobe
go build

## Prepull latest images
docker pull nginx
docker pull grpc/java-example-hostname

### HTTP Check Test
docker run --rm -p 8080:80 -d --name nginx-lprobe-test nginx 
echo "Wait 5s"
sleep 5
./lprobe -mode=http -port=8080 -endpoint=/
echo $?

if [ "$?" != 0 ]; then
    echo "HTTP test failed"
    docker stop nginx-lprobe-test
    exit 1
fi
docker stop nginx-lprobe-test


docker run --rm -p 8080:50051 -d --name grpc-lprobe-test grpc/java-example-hostname
echo "Wait 7s"
sleep 7
./lprobe -mode=grpc -port=8080 
echo $?
if [ "$?" != 0 ]; then
    echo "gRPC test failed"
    docker stop grpc-lprobe-test
    exit 1
fi
docker stop grpc-lprobe-test

echo "All good"
exit 0