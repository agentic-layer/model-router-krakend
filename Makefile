-include .env
export

VERSION=$(shell git describe --tags --always --first-parent)

.PHONY: help clean plugins test image up

help:
	@echo "Supported make targets (you can set the version in the Makefile):"
	@echo ""
	@echo "     clean   clean up build artifacts"
	@echo "   plugins   build and test all plugins"
	@echo "     image   build docker image and tag as latest and $(VERSION)"
	@echo "       run   build and start in local docker"
	@echo ""

.DEFAULT_GOAL := help

clean:
	rm -rf build

plugins:
	$(MAKE) -C ./go plugins

test:
	$(MAKE) -C ./go test

image:
	@echo VERSION=$(VERSION)
	docker build \
		--platform linux/amd64 \
		--tag model-router-krakend:$(VERSION) \
		.
	docker tag model-router-krakend:$(VERSION) model-router-krakend:latest

run: image
	docker run -p 8080:8080 -p 9090:9090 -e OPENAI_API_KEY=$(OPENAI_API_KEY) model-router-krakend:latest
