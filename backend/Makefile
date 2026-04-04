DIST_DIR=../infra/functions
ZIP_NAME=function.zip

FUNCTIONS=cmd/assets cmd/videos

.PHONY: all build zip clean

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

clean:
	@echo "Cleaning..."
	rm -rf $(DIST_DIR)


swag-init:
	@swag init -g cmd/assets/main.go -d . -o ./cmd/assets/docs --exclude "internal/modules/videos" --parseDependency --parseInternal
	@swag init -g cmd/videos/main.go -d . -o ./cmd/videos/docs --exclude "internal/modules/assets" --parseDependency --parseInternal