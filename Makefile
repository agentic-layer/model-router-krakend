VERSION=$(shell git describe --tags --always --first-parent)

.PHONY: help clean plugins test image up

help:
	@echo "Supported make targets (you can set the version in the Makefile):"
	@echo ""
	@echo "   plugins   build and test all plugins"
	@echo "     image   build docker image and tag as latest and $(VERSION)"
	@echo "        up   build and start in local docker"
	@echo ""

clean:
	rm -rf build

plugins:
	$(MAKE) -C ./go plugins

test:
	$(MAKE) -C ./go test

image:
	@echo VERSION=$(VERSION)
	docker build \
		--no-cache \
		--platform linux/amd64 \
		--build-arg GATEWAY_VERSION=$(VERSION) \
		--tag model-router-krakend:$(VERSION) \
		.
	docker tag model-router-krakend:$(VERSION) model-router-krakend:latest

up:
	docker compose up --build -d
