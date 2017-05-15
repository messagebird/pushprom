BINARY  = pushprom
PROJECT = pushprom

VERSION = $(shell git rev-list HEAD | wc -l |tr -d ' ')
HASH    = $(shell git rev-parse --short HEAD)

GO      = env GOPATH="$(PWD)/vendor:$(PWD)" go

all:
	@echo "make release       # Build $(PROJECT) for release"
	@echo "make development   # Build $(PROJECT) for development"
	@echo
	@echo "make run           # Run a development version of $(PROJECT)"
	@echo
	@echo "make container     # Create a Docker container for $(PROJECT)"
	@echo
	@echo "make test          # Run the test suite"
	@echo "make clean         # Clean up the project directory"


release: clean dependencies
	@echo "* Building $(PROJECT) for release"
	@$(GO) install $(PROJECT)/...

release_linux: export GOOS=linux
release_linux: export GOARCH=amd64
release_linux: release
	@mv bin/linux_amd64/${BINARY} bin/

development: clean dependencies
	@echo "* Building $(PROJECT) for development"
	@$(GO) install -race $(PROJECT)/...

run: development
	@echo "* Running development $(PROJECT) binary"
	@./bin/$(BINARY)

container: release_linux
	@echo "* Creating $(PROJECT) Docker container"
	@docker build -t $(PROJECT):$(VERSION) .
	@docker tag $(PROJECT):$(VERSION) $(PROJECT):latest

test:
	@echo "* Running tests"
	@$(GO) test $(PROJECT)/...

clean:
	rm -fr bin pkg vendor/pkg

dependencies:
	@echo "* go getting all dependencies into vendor/"
	@$(GO) get -t $(PROJECT)/...
	find vendor/ -name .git -type d | xargs rm -rf
