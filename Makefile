BINARY  = pushprom
PROJECT = pushprom

VERSION = $(shell git rev-list HEAD | wc -l |tr -d ' ')
HASH    = $(shell git rev-parse --short HEAD)

GO      = env GOPATH="$(PWD)/vendor:$(PWD)" go
LDFLAGS = -X main.BuildNumber=$(VERSION) -X main.CommitHash=$(HASH)

all:
	@echo "make release       # Build $(PROJECT) for release"
	@echo "make development   # Build $(PROJECT) for development"
	@echo
	@echo "make run           # Run a development version of $(BINARY)"
	@echo
	@echo "make container     # Create a Docker container for $(PROJECT)"
	@echo "make rollout       # Roll out the container to production"
	@echo
	@echo "make test          # Run the test suite"
	@echo "make clean         # Clean up the project directory"

release: clean
	@echo "* Building $(PROJECT) for release"
	@$(GO) install -ldflags '$(LDFLAGS)' $(PROJECT)/...

development: clean
	@echo "* Building $(PROJECT) for development"
	@$(GO) install -ldflags '$(LDFLAGS)' -race $(PROJECT)/...

dependencies:
	@echo "* go getting all dependencies into vendor/"
	@$(GO) get -t $(PROJECT)/...
	find vendor/ -name .git -type d | xargs rm -rf

run: development
	@echo "* Running development $(PROJECT) binary"
	@./bin/$(BINARY) -mapping-config=./mapping.conf -log.level=debug

test:
	@echo "* Running tests"
	@$(GO) test $(PROJECT)/...
	@echo

	@echo "* Checking code with golint"
	@golint src/$(PROJECT)/...
	@echo

	@echo "* Checking code with go vet"
	@$(GO) vet $(PROJECT)/...

clean:
	rm -fr bin pkg vendor/pkg lib/*.a lib/*.o
