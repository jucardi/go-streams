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

test-deps:
	@echo "installing test dependencies..."
	@go get github.com/smartystreets/goconvey/convey
	@go get gopkg.in/h2non/gock.v1
	@go get github.com/stretchr/testify/assert
	@go get github.com/axw/gocov/...
	@go get github.com/AlekSi/gocov-xml
	@go get gopkg.in/matm/v1/gocov-html

test: test-deps
	@echo "running test coverage..."
	@mkdir -p test-artifacts/coverage
	@gocov test ./... -v > test-artifacts/gocov.json
	@cat test-artifacts/gocov.json | gocov report
	@cat test-artifacts/gocov.json | gocov-xml > test-artifacts/coverage/coverage.xml
	@cat test-artifacts/gocov.json | gocov-html > test-artifacts/coverage/coverage.html

build: deps
	@echo "building..."
	@go build ./...
