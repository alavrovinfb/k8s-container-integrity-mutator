BINARY_NAME := mutator_app
APP_NAME := mutator
DEMO_NAME := demo
RELEASE_NAME_DB ?= db
POSTGRES_SERVICE ?= $(RELEASE_NAME_DB)-postgresql
NAMESPACE ?= default

IMAGE_NAME ?= integrity-injector
GIT_COMMIT := $(shell git describe --tags --long --dirty=-unsupported --always || echo pre-commit)
IMAGE_VERSION ?= $(GIT_COMMIT)

# helm chart path
HELM_CHART_PATH := helm-charts

## Downloads the Go module.
.PHONY : mod-download
mod-download:
	@echo "==> Downloading Go module"
	go mod download

## Downloads the necessesary dev dependencies.
.PHONY : dev-dependencies
dev-dependencies: minikube update docker helm-all
	@echo "==> Downloaded development dependencies"
	@echo "==> Successfully installed"

.PHONY : docker
docker:
	@eval $$(minikube docker-env) ;\
    docker build -t $(IMAGE_NAME):latest  -t $(IMAGE_NAME):$(IMAGE_VERSION) -f Dockerfile .

.PHONY: helm
helm-mutator:
	@helm upgrade -i ${APP_NAME} \
		--namespace=$(NAMESPACE) \
		--create-namespace \
		--set sideCar.secretName=$(POSTGRES_SERVICE) \
		--set sideCar.db.host=$(POSTGRES_SERVICE) \
		--set sideCar.db.name=$(DB_NAME) \
		--set sideCar.db.username=$(DB_USER) \
		--set sideCar.db.password=$(DB_PASSWORD) \
		--set image.repository=$(IMAGE_NAME)    \
		--set image.tag=$(IMAGE_VERSION) \
		$(HELM_CHART_PATH)/$(IMAGE_NAME)

helm-demo:
	@helm upgrade -i ${DEMO_NAME} \
		--namespace=$(NAMESPACE) \
		--create-namespace \
		$(HELM_CHART_PATH)/demo-app-to-inject

helm-demo-with-db:
	@helm upgrade -i ${DEMO_NAME} \
		--namespace=$(NAMESPACE) \
		--create-namespace \
		--set global.postgresql.auth.database=$(DB_NAME) \
		--set global.postgresql.auth.username=$(DB_USER) \
		--set global.postgresql.auth.password=$(DB_PASSWORD) \
		--set postgresql.enabled=true \
		--set postgresql.fullnameOverride=$(POSTGRES_SERVICE) \
		$(HELM_CHART_PATH)/demo-app-to-inject

.PHONY : tidy
tidy: ## Cleans the Go module.
	@echo "==> Tidying module"
	@go mod tidy

.PHONY : build
build:
	go build -o ${BINARY_NAME} cmd/main.go

.PHONY : run
run:
	go build -o ${BINARY_NAME} cmd/main.go
	./${BINARY_NAME}

## Cleaning build cache.
.PHONY : clean
clean:
	go clean
	rm ${BINARY_NAME}

## Compile the binary.
compile-all: windows-32bit windows-64bit linux-32bit linux-64bit MacOS

windows-32bit:
	echo "Building for Windows 32-bit"
	GOOS=windows GOARCH=386 go build -o bin/${BINARY_NAME}_32bit.exe cmd/main.go

windows-64bit:
	echo "Building for Windows 64-bit"
	GOOS=windows GOARCH=amd64 go build -o bin/${BINARY_NAME}_64bit.exe cmd/main.go

linux-32bit:
	echo "Building for Linux 32-bit"
	GOOS=linux GOARCH=386 go build -o bin/${BINARY_NAME}_32bit cmd/main.go

linux-64bit:
	echo "Building for Linux 64-bit"
	GOOS=linux GOARCH=amd64 go build -o bin/${BINARY_NAME} cmd/main.go

MacOS:
	echo "Building for MacOS X 64-bit"
	GOOS=darwin GOARCH=amd64 go build -o bin/${BINARY_NAME}_macos cmd/main.go
