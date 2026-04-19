DIST_DIR=../infra/artifacts/lambdas
ZIP_NAME=function.zip

FUNCTIONS=assets videos

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
	@swag init -g cmd/assets/main.go -d . -o ./cmd/assets/docs --exclude "internal/modules/campaigns" --parseDependency --parseInternal
	@swag init -g cmd/campaigns/main.go -d . -o ./cmd/campaigns/docs --exclude "internal/modules/assets" --parseDependency --parseInternal