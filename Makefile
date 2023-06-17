.PHONY: test format vet staticcheck magneticod magneticow

all: test magneticod magneticow

magneticod:
	go install --tags fts5 ./cmd/magneticod

magneticow:
	# https://github.com/kevinburke/go-bindata
	go-bindata -pkg "main" -o="cmd/magneticow/bindata.go" -prefix="cmd/magneticow/data/" cmd/magneticow/data/...
	# Prepend the linter instruction to the beginning of the file
	#sed -i '1s;^;//lint:file-ignore * Ignore file altogether\n;' cmd/magneticow/bindata.go
	go install --tags fts5 ./cmd/magneticow

vet:
	go vet ./...

test:
	go test ./...

format:
	gofmt -w ./cmd/
	gofmt -w ./pkg/
