test:
	go test -v ./...
cli:
	go build -o migrate-passwd migrate-passwd.go
lint:
	find . -type f -name '*.go' | xargs goimports -w
	find . -type f -name '*.go' | xargs -n 1 go fmt
