build:
	go build -v -o bin/boa ./cmd/boa

install:
	go install -v ./cmd/boa

test:
	go test -v ./...

golden:
	go test -v $$(grep -lR --include=\*_test.go "flag.Bool(\"update\"" . | xargs dirname  | sort  | uniq ) -update

lint:
	golangci-lint run
