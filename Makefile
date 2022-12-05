BINARY  = pushprom
PROJECT = pushprom

VERSION = $(shell git rev-list HEAD | wc -l |tr -d ' ')
HASH    = $(shell git rev-parse --short HEAD)

GO      = go

all:
	@echo "make container     # Create a Docker container for $(PROJECT)"
	@echo
	@echo "make test          # Run the test suite"


release_linux: 
	go mod download
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "-s -w" -o "bin/pushprom-$(VERSION).linux-amd64/pushprom" github.com/messagebird/pushprom
	mv bin/pushprom-$(VERSION).linux-amd64/pushprom bin/
	rm -rf bin/pushprom-$(VERSION).linux-amd64/
	file bin/pushprom

native: 
	go clean
	go mod download
	CGO_ENABLED=0 go build -ldflags "-s -w" -o "bin/pushprom-native" github.com/messagebird/pushprom
	file bin/pushprom-native

container: 
	echo "* Creating $(PROJECT) Docker container"
	DOCKER_BUILDKIT=1 docker build -t $(PROJECT):$(VERSION) .
	docker tag $(PROJECT):$(VERSION) $(PROJECT):latest
