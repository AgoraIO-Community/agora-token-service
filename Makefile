SHELL := /bin/bash

TAG_SHA = agora-token-service

# set the build args from env file
DECONARGS = $(shell echo "$$(for i in `cat .env`; do out+="--build-arg $$i " ; done; echo $$out;out="")")
# Generate the Arguments using DECONARGS
GEN_ARGS = $(eval BARGS=$(DECONARGS))
# Set the SERVER_PORT from .env, run target checks if set and defaults to 8080 if needed
SERVER_PORT = $(shell grep SERVER_PORT .env | cut -d '=' -f2 | tr -d '[:space:]' || echo "8080")

# Docker and Golang source files
DOCKER_FILES := Dockerfile
GO_SOURCE_FILEs := $(shell find . -type f -name '*.go')


.PHONY: all check-env build run clean

all: build run

check-env: 
	@if [ ! -f .env ]; then \
		echo ".env file not found. Please create one."; \
		exit 1;\
	fi

build_marker: $(DOCKER_FILES) $(GO_SOURCE_FILEs) .env
	@echo "Running docker build with tag: ${TAG_SHA}"
		$(GEN_ARGS)
	docker build -t $(TAG_SHA) $(BARGS) .
	@touch build_marker

build: check-env build_marker

run: 
	@SERVER_PORT=$${SERVER_PORT:-8080}; \
	echo "Running docker container on port: $$SERVER_PORT"; \
	docker run -p $$SERVER_PORT:$$SERVER_PORT agora-token-service

clean:
	rm -f build_marker