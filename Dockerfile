FROM golang:1.10 as builder
WORKDIR /go/src/github.com/vigo5190/go-socks5
COPY . .
RUN CGO_ENABLED=0 GOOS=linux make

FROM alpine:latest
WORKDIR /root/

COPY --from=builder /go/src/github.com/vigo5190/go-socks5 .

EXPOSE 8008
ENTRYPOINT [ "./go-socks5"]

