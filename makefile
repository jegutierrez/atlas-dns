lint:
	@echo "running linter"
	@golangci-lint run ./...

test:
	@echo "running tests"
	@go test ./... -v -covermode=atomic -coverpkg=./... -count=1 -race