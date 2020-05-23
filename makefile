XC_OS="linux darwin"
XC_ARCH="amd64"
XC_PARALLEL="2"
XC_LDFLAGS="-s -w"
BIN="./bin"
GOGET=go get
SRC=$(shell find . -name "*.go")

ifeq (, $(shell which gox))
$(warning "could not find gox in $(PATH), run: go get github.com/mitchellh/gox")
endif

.PHONY: all deps build

default: all

all: deps build

deps:
		$(GOGET) github.com/spf13/cobra
		$(GOGET) github.com/aws/aws-sdk-go/aws
		$(GOGET) github.com/aws/aws-sdk-go/aws/session
		$(GOGET) github.com/aws/aws-sdk-go/service/ecs
		$(GOGET) github.com/aws/aws-sdk-go/service/elbv2
		$(GOGET) github.com/jedib0t/go-pretty/table
build:
	gox \
		-os=$(XC_OS) \
		-arch=$(XC_ARCH) \
		-parallel=$(XC_PARALLEL) \
		-ldflags=$(XC_LDFLAGS) \
		-output=$(BIN)/{{.Dir}}_{{.OS}}_{{.Arch}} \
		;