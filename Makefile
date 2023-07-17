.PHONY: test format vet staticcheck magneticod magneticow

all: test magneticod magneticow

magneticod:
	go install --tags fts5 ./cmd/magneticod

magneticow:
	go install --tags fts5 ./cmd/magneticow

vet:
	go vet ./...

test:
	go test ./...

format:
	gofmt -w ./cmd/
	gofmt -w ./pkg/
