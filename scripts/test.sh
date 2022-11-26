#!/usr/bin/env bash

set -e
rm -rf ./lprobe
go build
docker pull nginx
docker run --rm -p 8080:80 -d --name nginx-lprobe-test nginx 
sleep 5
./lprobe -port=8080 -endpoint=/
echo $?

if [ "$?" != 0 ]; then
    echo "test failed"
    docker stop nginx-lprobe-test
    exit 1
fi
set -e

echo "All good"
docker stop nginx-lprobe-test
exit 0