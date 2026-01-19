.PHONY: check format test

format:
	@go tool goimports -w .

check: format
	@go vet
	@go tool staticcheck

test: check
	@go test
