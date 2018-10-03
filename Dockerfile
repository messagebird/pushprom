# ##### builder stage #####
FROM alpine:latest AS builder
LABEL maintainer="support@messagebird.com"

ENV GOPATH=/usr/local

RUN apk add --no-cache musl-dev go git && \
    go get -u github.com/golang/dep/cmd/dep

ADD . /usr/local/src/pushprom
WORKDIR /usr/local/src/pushprom

RUN dep ensure -v && \
    go test && \
    GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build -ldflags "-s -w" -o "/usr/local/bin/pushprom"

# ##### image stage #####
FROM alpine:latest
COPY --from=builder /usr/local/bin/pushprom /usr/local/bin/pushprom

EXPOSE 9090 9091

# Run
CMD ["/usr/local/bin/pushprom"]

