FROM golang:alpine AS builder
RUN mkdir /build
ADD . /build/
WORKDIR /build
RUN go build -o check

FROM alpine:latest
RUN apk --no-cache add aspell aspell-en
COPY --from=builder /build/check /check
WORKDIR /
ENTRYPOINT ["/check"]
