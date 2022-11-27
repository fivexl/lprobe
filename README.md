# lprobe
A command-line tool to perform Local Health Check Probes inside Container Images (ECS, Docker)
## HOW TO
### Local run
```shell
./lprobe -port=8080 -endpoint=/
```

### Add to a container image
You can bundle the statically compiled lprobe in your container image. Choose a binary release and download it in your Dockerfile:
```
RUN LPROBE_VERSION=v0.0.1 && \
    wget -qO/bin/lprobe https://github.com/fivexl/lprobe/releases/download/${LPROBE_VERSION}/lprobe-linux-amd64 && \
    chmod +x /bin/lprobe
```

### Docker Healthcheck 
```
HEALTHCHECK  --interval=5m --timeout=3s \
  CMD lprobe -port=8080 -endpoint=/healthz
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

# Source code used
- https://github.com/grpc-ecosystem/grpc-health-probe
