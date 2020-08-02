NAME=dns

lint:
	@echo "running linter."
	@golangci-lint run ./... -v

test:
	@echo "running tests."
	@go test ./... -v -covermode=atomic -coverpkg=./... -count=1 -race

build:
	go build -o ${NAME} "./cmd/server"

clean:
	rm -f ${NAME}

run: build
	DNS_PORT=8080 DNS_SECTOR_ID=1 DNS_LOG_LEVEL=DEBUG ./dns

docker-build:
	@echo "building DNS docker image."
	docker build -t jegutierrez/atlas-dns:latest .