.PHONY: help build build-python build-go build-distroless build-distroless-python build-distroless-go \
        push push-python push-go pull pull-python pull-go \
        run-python run-go run-distroless-python run-distroless-go \
        compose-up compose-down compose-build compare-sizes clean

# Docker Hub username
DOCKER_USER = escalopax

# Image names and tags
PYTHON_IMAGE = $(DOCKER_USER)/moscow-time-app
GO_IMAGE = $(DOCKER_USER)/wordle-game
TAG = latest

# Distroless image tags
PYTHON_DISTROLESS_IMAGE = $(DOCKER_USER)/moscow-time-app-distroless
GO_DISTROLESS_IMAGE = $(DOCKER_USER)/wordle-game-distroless

help: ## Show this help message
	@echo "Usage: make [target]"
	@echo ""
	@echo "Available targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-25s\033[0m %s\n", $$1, $$2}'

## Build targets
build: build-python build-go ## Build both Python and Go applications

build-python: ## Build Python application Docker image
	@echo "Building Python application..."
	docker build -t $(PYTHON_IMAGE):$(TAG) ./app_python

build-go: ## Build Go application Docker image
	@echo "Building Go application..."
	docker build -t $(GO_IMAGE):$(TAG) ./app_go

build-distroless: build-distroless-python build-distroless-go ## Build both distroless images

build-distroless-python: ## Build Python distroless Docker image
	@echo "Building Python distroless application..."
	docker build -f ./app_python/distroless.Dockerfile -t $(PYTHON_DISTROLESS_IMAGE):$(TAG) ./app_python

build-distroless-go: ## Build Go distroless Docker image
	@echo "Building Go distroless application..."
	docker build -f ./app_go/distroless.Dockerfile -t $(GO_DISTROLESS_IMAGE):$(TAG) ./app_go

## Push targets
push: push-python push-go ## Push both images to Docker Hub

push-python: ## Push Python image to Docker Hub
	@echo "Pushing Python image to Docker Hub..."
	docker push $(PYTHON_IMAGE):$(TAG)

push-go: ## Push Go image to Docker Hub
	@echo "Pushing Go image to Docker Hub..."
	docker push $(GO_IMAGE):$(TAG)

push-distroless: ## Push distroless images to Docker Hub
	@echo "Pushing distroless images to Docker Hub..."
	docker push $(PYTHON_DISTROLESS_IMAGE):$(TAG)
	docker push $(GO_DISTROLESS_IMAGE):$(TAG)

## Pull targets
pull: pull-python pull-go ## Pull both images from Docker Hub

pull-python: ## Pull Python image from Docker Hub
	@echo "Pulling Python image from Docker Hub..."
	docker pull $(PYTHON_IMAGE):$(TAG)

pull-go: ## Pull Go image from Docker Hub
	@echo "Pulling Go image from Docker Hub..."
	docker pull $(GO_IMAGE):$(TAG)

## Run targets
run-python: ## Run Python application container
	@echo "Running Python application..."
	docker run -d -p 5000:5000 --name moscow-time-app $(PYTHON_IMAGE):$(TAG)

run-go: ## Run Go application container
	@echo "Running Go application..."
	docker run -d -p 8080:8080 --name wordle-game-app $(GO_IMAGE):$(TAG)

run-distroless-python: ## Run Python distroless container
	@echo "Running Python distroless application..."
	docker run -d -p 5000:5000 --name moscow-time-app-distroless $(PYTHON_DISTROLESS_IMAGE):$(TAG)

run-distroless-go: ## Run Go distroless container
	@echo "Running Go distroless application..."
	docker run -d -p 8080:8080 --name wordle-game-app-distroless $(GO_DISTROLESS_IMAGE):$(TAG)

## Docker Compose targets
compose-up: ## Start applications using Docker Compose (builds from source)
	@echo "Starting applications with Docker Compose..."
	docker compose up -d

compose-down: ## Stop applications using Docker Compose
	@echo "Stopping applications with Docker Compose..."
	docker compose down

compose-build: ## Build and start applications using Docker Compose
	@echo "Building and starting applications with Docker Compose..."
	docker compose up --build -d

compose-logs: ## Show Docker Compose logs
	docker compose logs -f

compose-prod-up: ## Start applications using production images from Docker Hub
	@echo "Starting applications with production images from Docker Hub..."
	docker compose -f docker-compose.prod.yml up -d

compose-prod-down: ## Stop applications using production compose
	@echo "Stopping applications..."
	docker compose -f docker-compose.prod.yml down

compose-prod-pull: ## Pull latest production images from Docker Hub
	@echo "Pulling latest images from Docker Hub..."
	docker compose -f docker-compose.prod.yml pull

## Utility targets
compare-sizes: ## Compare image sizes (regular vs distroless)
	@echo "============================================"
	@echo "Image Sizes Comparison"
	@echo "============================================"
	@echo "\nCommand: docker images | grep escalopax"
	@echo "============================================"
	@echo "\nRegular Images:"
	@docker images | grep -E "$(PYTHON_IMAGE)|$(GO_IMAGE)" | grep -v distroless || echo "Regular images not found. Run 'make build' first."
	@echo "\nDistroless Images:"
	@docker images | grep -E "distroless" | grep -E "$(DOCKER_USER)" || echo "Distroless images not found. Run 'make build-distroless' first."
	@echo "============================================"

stop: ## Stop all running containers
	@echo "Stopping all running containers..."
	-docker stop moscow-time-app wordle-game-app moscow-time-app-distroless wordle-game-app-distroless 2>/dev/null || true

clean: stop ## Clean up containers and images
	@echo "Cleaning up containers..."
	-docker rm moscow-time-app wordle-game-app moscow-time-app-distroless wordle-game-app-distroless 2>/dev/null || true
	@echo "Cleaning up images..."
	-docker rmi $(PYTHON_IMAGE):$(TAG) $(GO_IMAGE):$(TAG) 2>/dev/null || true
	-docker rmi $(PYTHON_DISTROLESS_IMAGE):$(TAG) $(GO_DISTROLESS_IMAGE):$(TAG) 2>/dev/null || true

lint-dockerfile: ## Lint Dockerfiles using hadolint
	@echo "Linting Dockerfiles..."
	@command -v hadolint >/dev/null 2>&1 || { echo "hadolint not installed. Install from https://github.com/hadolint/hadolint"; exit 1; }
	@hadolint app_python/Dockerfile
	@hadolint app_go/Dockerfile
	@hadolint app_python/distroless.Dockerfile
	@hadolint app_go/distroless.Dockerfile

test-running: ## Test applications are running
	@echo "Testing Python application..."
	@curl -f http://localhost:5000/health || echo "Python app not responding"
	@echo "\nTesting Go application..."
	@curl -f http://localhost:8080/health || echo "Go app not responding"

test-python: ## Run Python unit tests with coverage
	@echo "Running Python tests with coverage..."
	cd app_python && pytest --cov=app --cov-report=term --cov-report=term-missing
	@echo ""
	@echo "📊 Coverage Summary:"
	cd app_python && coverage report --format=total

test-go: ## Run Go unit tests with coverage
	@echo "Running Go tests with coverage..."
	cd app_go && go test -v -race -coverprofile=coverage.out -covermode=atomic ./...
	@echo ""
	@echo "📊 Coverage Summary:"
	cd app_go && go tool cover -func=coverage.out | tail -1

test-all: test-python test-go ## Run all unit tests

lint-python: ## Lint Python code
	@echo "Linting Python code..."
	cd app_python && pylint app.py --rcfile=.pylintrc

lint-go: ## Lint Go code
	@echo "Linting Go code..."
	cd app_go && revive -config .revive.toml -formatter friendly ./...

lint-all: lint-python lint-go ## Run all linters

ci-local: lint-all test-all ## Run CI checks locally
