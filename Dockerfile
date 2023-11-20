FROM golang:1.21 as builder
WORKDIR /go/src/github.com/szkrstf/packs
COPY . .
RUN make build-linux

FROM alpine
RUN apk add --no-cache ca-certificates
COPY --from=builder /go/src/github.com/szkrstf/packs/packs /app
EXPOSE 8080
CMD ["/app"]