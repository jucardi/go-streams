all: format vet deps test build

format:
	@echo "formatting files..."
	@go get golang.org/x/tools/cmd/goimports
	@goimports -w -l .
	@gofmt -s -w -l .

vet:
	@echo "vetting..."
	@go vet ./...

deps:
	@echo "installing dependencies..."
	@go get ./...

test:
	@echo "running test coverage..."
	@mkdir -p test-artifacts/coverage
	@go test -p 1 ./... -v -coverprofile test-artifacts/cover.out
	@go tool cover -func test-artifacts/cover.out

build: deps
	@echo "building..."
	@go build ./...
