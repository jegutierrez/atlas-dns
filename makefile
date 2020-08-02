NAME=dns

lint:
	@echo "running linter."
	@golangci-lint run ./... -v

test:
	@echo "running tests."
	@go test ./... -v -covermode=atomic -coverpkg=./... -count=1 -race

docker-build:
	@echo "building DNS docker image."
	docker build -t jegutierrez/atlas-dns:latest .