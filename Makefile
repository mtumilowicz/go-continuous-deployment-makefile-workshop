# Makefile

# Variables
# COMMIT_HASH ?= f7849357b169da3d0f446b2717a3eef644159fdd
COMMIT_HASH ?= 292e275cacf0238ec0e3d76e8c4948a02c051fc7
IMAGE_NAME ?= helm-workshop
RELEASE_NAME ?= helmworkshopchart
REPO_URL ?= https://github.com/mtumilowicz/helm-workshop
CHART_DIR ?= ./helm
NAMESPACE ?= default

# Build the Go binary
build:
	go build -o bin/deployer main.go

# Target: clone repository
clone: build
	./deployer -action clone -repo-url $(REPO_URL) -image-name $(IMAGE_NAME)

# Target: checkout commit
checkout: build
	./deployer -action checkout -commit-hash $(COMMIT_HASH) -image-name $(IMAGE_NAME)

# Target: clean build artifacts
clean: build
	./deployer -action clean -image-name $(IMAGE_NAME)

# Target: run tests
test: build
	./deployer -action test -image-name $(IMAGE_NAME)

# Target: build Docker image
build-image: build
	./deployer -action build -image-version $(COMMIT_HASH) -image-name $(IMAGE_NAME)

# Target: upgrade Helm chart
upgrade: build
	./deployer -action upgrade -image-version $(COMMIT_HASH) -image-name $(IMAGE_NAME) -release-name $(RELEASE_NAME) -chart-dir $(CHART_DIR) -namespace $(NAMESPACE)

# Phony targets
.PHONY: build clone checkout clean test build-image upgrade
