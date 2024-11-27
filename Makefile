test:
	@echo "Running tests"
	@go test -count=1 ./...

run:
	@echo "Running application"
	@go run cmd/main.go

build:
	@echo "Building application"
	@go build -o bin/client-server-api cmd/main.go

clean:
	@echo "Cleaning up"
	@rm -rf bin

fmt:
	@echo "Formatting code"
	@go fmt ./...