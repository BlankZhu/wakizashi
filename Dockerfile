# docker build -t probe:1.0.0 -f ./Dockerfile --build-arg APP_NAME=appname VERSION=1.0.0 .

FROM golang:1.14-alpine3.12 as builder

RUN apk add --no-cache git &&\
    apk add --no-cache build-base &&\
    apk add --no-cache linux-headers

# APP_NAME=probe, center
ARG APP_NAME
ARG VERSION

COPY . .
WORKDIR /usr/wakizashi
RUN GO111MODULE=on go build -ldflags \
    "-X main.buildTime=`date +%Y-%m-%d,%H:%M:%S` -X main.buildVersion=${VERSION} -X main.gitCommitID=`git rev-parse HEAD`" \
    -o main \
    ./cmd/${APP_NAME} && \
    mv /usr/wakizashi/config/${APP_NAME}-config.yaml /usr/wakizashi/config/config.yaml
##

FROM alpine:latest

RUN mkdir -p /usr/wakizashi &&\
    mkdir -p /usr/wakizashi/config &&

WORKDIR /usr/wakizashi

COPY --from=builder /usr/wakizashi/main /usr/wakizashi/
COPY --from=builder /usr/wakizashi/config/config.yaml /usr/wakizashi/config

ENTRYPOINT [ "./main", "-c", "./config/config.yaml", "-v", "true"]