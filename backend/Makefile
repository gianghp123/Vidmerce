DIST_DIR=../infra/artifacts/lambdas
ZIP_NAME=function.zip

FUNCTIONS=assets webhooks

.PHONY: all build zip swag-init

all: build zip

build:
	@echo "Building Go Lambdas..."
	@for fn in $(FUNCTIONS); do \
		echo "Building $$fn..."; \
		mkdir -p $(DIST_DIR)/$$fn; \
		GOOS=linux GOARCH=amd64 CGO_ENABLED=0 \
		go build -o $(DIST_DIR)/$$fn/bootstrap ./$$fn/main.go; \
	done

zip:
	@echo "Zipping Lambda packages..."
	@for fn in $(FUNCTIONS); do \
		echo "Zipping $$fn..."; \
		cd $(DIST_DIR)/$$fn && zip -q $(ZIP_NAME) bootstrap && cd - > /dev/null; \
	done


swag-init:
	@swag init \
		-g main.go \
		-d "./cmd/assets,internal/modules/assets" \
		-o ./cmd/assets/docs \
		--parseDependency \
		--parseInternal

	@swag init \
		-g main.go \
		-d "./cmd/webhooks,internal/modules/webhooks" \
		-o ./cmd/webhooks/docs \
		--parseDependency \
		--parseInternal

	@swag init \
		-g main.go \
		-d "./cmd/campaigns,internal/modules/campaigns" \
		-o ./cmd/campaigns/docs \
		--parseDependency \
		--parseInternal