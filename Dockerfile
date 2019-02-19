FROM golang:alpine as build-env
# All these steps will be cached

RUN apk add git
RUN mkdir /pushprom
WORKDIR /pushprom
COPY go.mod . 
COPY go.sum .

# Get dependencies - will also be cached if we won't change mod/sum
RUN go mod download
# COPY the source code as the last step
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -ldflags "-s -w" -o /go/bin/pushprom .

# <- Second step to build minimal image
FROM scratch 
COPY --from=build-env /go/bin/pushprom /go/bin/pushprom
ENTRYPOINT ["/go/bin/pushprom"]
