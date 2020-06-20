build:
	go build -v -o bin/boa ./cmd/boa

test:
	go test -v ./...

lint:
	golangci-lint run

