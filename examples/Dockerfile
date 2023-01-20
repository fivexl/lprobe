FROM nginx:alpine
ARG LPROBE_VERSION=v0.0.5
ARG TARGETPLATFORM
RUN case ${TARGETPLATFORM} in \
         "linux/amd64")  LPROBE_ARCH=amd64  ;; \
         "linux/arm64")  LPROBE_ARCH=arm64  ;; \
    esac \
 && wget -qO/bin/lprobe https://github.com/fivexl/lprobe/releases/download/${LPROBE_VERSION}/lprobe-linux-${LPROBE_ARCH} \
 && chmod +x /bin/lprobe \
 && rm -f /usr/bin/wget
HEALTHCHECK --interval=15s --timeout=5s --start-period=5s --retries=3 CMD [ "lprobe", "-mode=http", "-port=80", "-endpoint=/" ]
