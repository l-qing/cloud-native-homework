FROM golang:1.17 as builder

ADD . /workspace
WORKDIR /workspace

ENV GO111MODULE=on \
    GOPROXY=goproxy.cn,direct \
    CGO_ENABLED=0 \
    GOOS=linux

RUN go build -o httpserver .


FROM alpine:3.14

ARG version=v1.0.0
ENV VERSION=${version}
EXPOSE 8081

COPY --from=builder /workspace/httpserver /httpserver
ENTRYPOINT ["/httpserver"]

