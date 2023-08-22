[![FivexL](https://releases.fivexl.io/fivexlbannergit.jpg)](https://fivexl.io/)

# Why LProbe?
A command-line tool to perform Local Health Check Probes inside Container Images (ECS, Docker, Kubernetes). When your container gets breached, the intruder/attacker can use tools like wget or curl to download more tools for further exploitation and lateral movement within your system. Thus we developed LProbe as wget/curl replacement for hardened and secure container images.

## HOW TO
### Local run
```shell
./lprobe -port=8080 -endpoint=/
```

### Local run for gRPC
```shell
./lprobe -port=8080 -mode=grpc
```

### Add to a container image
You can bundle the statically compiled LProbe in your container image. Choose a binary release and download it in your Dockerfile:
```
ARG LPROBE_VERSION=v0.0.6
ARG TARGETPLATFORM
RUN case ${TARGETPLATFORM} in \
         "linux/amd64")  LPROBE_ARCH=amd64  ;; \
         "linux/arm64")  LPROBE_ARCH=arm64  ;; \
    esac \
 && wget -qO/bin/lprobe https://github.com/fivexl/lprobe/releases/download/${LPROBE_VERSION}/lprobe-linux-${LPROBE_ARCH} \
 && chmod +x /bin/lprobe \
 && rm -f /usr/bin/wget
```

### Add to a container image from lprobe containter
```
FROM scratch
COPY --from=ghcr.io/fivexl/lprobe:0.0.6 /lprobe /bin/lprobe
```

### Docker Healthcheck 
```
HEALTHCHECK --interval=15s --timeout=5s --start-period=5s --retries=3 CMD [ "lprobe", "-mode=http", "-port=8080", "-endpoint=/healthz" ]
```

### ECS Healthcheck 
```
[ "CMD", "lprobe", "-port=8080", "-endpoint=/healthz"]
```

### Kubernetes (k8S) Healthcheck
```
spec:
  containers:
  - name: server
    image: "[YOUR-DOCKER-IMAGE]"
    ports:
    - containerPort: 8080
    readinessProbe:
      exec:
        command: ["/bin/lprobe", "-port=8080", "-endpoint=/readiness"]
      initialDelaySeconds: 5
    livenessProbe:
      exec:
        command: ["/bin/lprobe", "-port=8080", "-endpoint=/liveness"]
      initialDelaySeconds: 10
```

# Dev Guide
```
export GO111MODULE=on
go mod init lprobe
go mod tidy
go run .
```

```
go mod edit -go=1.21
go get -u -t -x ./...
go mod tidy
```

# Source code used
- https://github.com/grpc-ecosystem/grpc-health-probe
