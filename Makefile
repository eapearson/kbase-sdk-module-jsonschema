.PHONY: test compile

test:
	@echo "Testing 123..."

start:
	docker run .

compile:
	go build server.go

dev-image:
	@bash local-scripts/build-image.sh

run-dev-image:
	@bash local-scripts/run-image.sh