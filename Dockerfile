FROM golang:1.20 as builder

WORKDIR /go/src/
COPY . .

RUN export GOPROXY=https://goproxy.cn && \
    go build -o pizza-apiserver

FROM centos:7

COPY --from=builder /go/src/pizza-apiserver /pizza-apiserver

EXPOSE 8081

CMD ["/pizza-apiserver"]
