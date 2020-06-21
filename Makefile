build:
	go build -v -o bin/boa ./cmd/boa

install:
	go install -v ./cmd/boa

test:
	go test -v ./...

golden:
	go test -v ./... -update

lint:
	golangci-lint run
