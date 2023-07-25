FROM golang:alpine as build-env

# Copy source + vendor
COPY . /go/src/github.com/KscSDK/kingsoftcloud-exporter
WORKDIR /go/src/github.com/KscSDK/kingsoftcloud-exporter

# Build
ENV GOPATH=/go
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -v -a -ldflags "-s -w" -o /go/bin/ksc-exporter ./cmd/ksc-exporter/

FROM alpine
COPY --from=build-env /go/bin/ksc_exporter /usr/bin/ksc_exporter
RUN apk update
#RUN apk add git
RUN apk add curl
RUN apk add tcpdump
ENTRYPOINT ["ksyun_monitor_exporter"]
