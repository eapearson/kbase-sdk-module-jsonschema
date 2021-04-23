.PHONY: test compile

test:
	@echo "Testing 123..."

start:
	docker run .

compile:
	go build server.go

build-image:
	@bash scripts/build-image.sh

run-image:
	@bash scripts/run-image.sh