.PHONY: check format test

format:
	@go tool goimports -w .
	@go tool gofumpt -l -w -extra .

check: format
	@go vet
	@go tool staticcheck

test: check
	@go test
