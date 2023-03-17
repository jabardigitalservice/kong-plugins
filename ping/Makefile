.PHONY: lint
lint:
	@golangci-lint run -v

.PHONY: validate
validate:
	@go test -v ./...
	@golangci-lint run -v
	@go build -v ./...

.PHONY: build
build:
	@GOOS=linux GOARCH=amd64 go build -o ./bin/ping ./src
	
.PHONY: cover
cover:
	@go test -v -coverprofile cover.out ./...
	@go tool cover -html cover.out -o cover.html
