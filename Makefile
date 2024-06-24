build:
	@go build -o bin/web-server cmd/main.go

test:
	@go test -v ./...

run: build
	@./bin/web-server