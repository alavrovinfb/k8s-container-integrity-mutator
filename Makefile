BINARY_NAME=mutator_app
APP_NAME=mutator
DEMO_NAME=demo
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
    docker build -t mutator:latest -f Dockerfile .

.PHONY: helm
helm-mutator:
	helm install ${APP_NAME} helm-charts/mutator

helm-demo:
	helm install ${DEMO_NAME} helm-charts/demo-app-to-inject

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
