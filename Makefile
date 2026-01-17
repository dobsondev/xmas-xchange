build:
	go build -o bin/xmas-xchange main.go

run:
	go run .

debug:
	LOG_LEVEL=DEBUG go run .

test:
	go test ./...

test-verbose:
	go test -v ./...

test-coverage:
	go test -cover ./...

test-coverage-report:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"
	@echo "Link: file://$(PWD)/coverage.html"