#############      builder       #############
FROM golang:1.13.9 AS builder

ARG TARGETS=dev

WORKDIR /go/src/github.com/mandelsoft/metal-lookup
COPY . .

RUN make $TARGETS

############# base
FROM alpine:3.11.3 AS base

#############      main     #############
FROM base AS main
COPY --from=builder /go/bin/server /metal-lookup-server

WORKDIR /

ENTRYPOINT ["/metal-lookup-server"]
