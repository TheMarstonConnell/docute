install:
	@ go install .

generate:
	cd docs && docute generate

host: install generate
	cd docs && docute host


.PHONY: host install generate