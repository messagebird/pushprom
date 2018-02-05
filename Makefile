BINARY  = pushprom
PROJECT = pushprom

VERSION = $(shell git rev-list HEAD | wc -l |tr -d ' ')
HASH    = $(shell git rev-parse --short HEAD)

GO      = go

all:
	@echo "make container     # Create a Docker container for $(PROJECT)"
	@echo
	@echo "make test          # Run the test suite"


release_linux: export GOOS=linux
release_linux: export GOARCH=amd64
release_linux: 
	dep ensure
	go build -ldflags "-s -w" -o "bin/pushprom-$(VERSION).linux-amd64/pushprom" github.com/messagebird/pushprom	
	mv bin/pushprom-$(VERSION).linux-amd64/pushprom bin/

container: release_linux
	@echo "* Creating $(PROJECT) Docker container"
	@docker build -t $(PROJECT):$(VERSION) .
	@docker tag $(PROJECT):$(VERSION) $(PROJECT):latest
