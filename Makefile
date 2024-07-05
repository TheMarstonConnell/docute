install:
	@ go install .

generate:
	cd docs && docute generate

host: install generate
	cd docs && docute host



format-tools:
	go install mvdan.cc/gofumpt@v0.6.0
	gofumpt -l -w .

lint: format-tools
	golangci-lint run

.PHONY: host install generate format-tools lint