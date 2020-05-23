    # Go parameters
    GOCMD=go
    GOBUILD=$(GOCMD) build
    GOCLEAN=$(GOCMD) clean
    GOTEST=$(GOCMD) test
    GOGET=$(GOCMD) get
    BINARY_NAME=ecs-mon
    BINARY_UNIX=$(BINARY_NAME)_unix
    BINARY_DIR=bin
    
    all: test build
    build:
			mkdir -p $(BINARY_DIR) 
			$(GOBUILD) -o $(BINARY_NAME) -v
    test: 
			$(GOTEST) -v ./...
    clean: 
			$(GOCLEAN)
			rm -f $(BINARY_DIR)/*
    run:
			$(GOBUILD) -o $(BINARY_NAME) -v ./...
			./$(BINARY_NAME)
    deps:
			$(GOGET) github.com/spf13/cobra
			$(GOGET) github.com/aws/aws-sdk-go/aws
			$(GOGET) github.com/aws/aws-sdk-go/aws/session
			$(GOGET) github.com/aws/aws-sdk-go/service/ecs
			$(GOGET) github.com/aws/aws-sdk-go/service/elbv2
			$(GOGET) github.com/jedib0t/go-pretty/table
    
    # Cross compilation
    build:
		CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -ldflags="-s -w" -o $(BINARY_DIR)/$(BINARY_UNIX) -v
		CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GOBUILD) -ldflags="-s -w" -o $(BINARY_DIR)/$(BINARY_NAME) -v
		upx --brute $(BINARY_DIR)/$(BINARY_NAME)