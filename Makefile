build:
	@go build -o bin/proxy_server cmd/main.go

run: build
	@./bin/proxy_server