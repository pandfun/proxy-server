build:
	@go build -o bin/proxy-server cmd/main.go

test:
	@go test -v ./...

run: build
	@./bin/proxy-server