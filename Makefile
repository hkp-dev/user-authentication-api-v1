BINARY_NAME = app
GOLANG = go

GOARCH := $(shell uname -m)
ifeq ($(GOARCH),x86_64)
	GOARCH=amd64
else ifeq ($(GOARCH),aarch64)
	GOARCH=arm64
else ifeq ($(GOARCH),armv7l)
	GOARCH=arm
endif

GOOS := $(shell uname -s | tr A-Z a-z)
ifeq ($(GOOS),linux)
	GOOS=linux
else ifeq ($(GOOS),darwin)
	GOOS=darwin
else ifeq ($(GOOS),windows)
	GOOS = windows
endif
build:
	$(GOLANG) build -o $(BINARY_NAME)
	@echo("$(BINARY_NAME) built successfully for GOOS=$(GOOS) GOARCH=$(GOARCH)")
clean:
	rm -rf $(BINARY_NAME)
run:
	$(GOLANG) run .
