test:
	go test ./...

test-coverage:
	go test -race -covermode atomic -coverprofile=coverage.out ./...

coverage:
	go tool cover -func=coverage.out