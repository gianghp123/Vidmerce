# Biến mặc định
ENV ?= dev
TERRAFORM_DIR = infra/envs/$(ENV)
BACKEND_DIR = backend

.PHONY: build deploy clean

# Target build: Chạy makefile của backend
build:
  @echo "Building backend for environment: $(ENV)..."
  cd $(BACKEND_DIR) && $(MAKE)

# Target deploy: Build trước, sau đó apply terraform
deploy: build
  @echo "Deploying to environment: $(ENV) (Path: $(TERRAFORM_DIR))..."
  @if [ -d "$(TERRAFORM_DIR)" ]; then \
    cd $(TERRAFORM_DIR) && terraform init && terraform apply -auto-approve; \
  else \
    echo "Error: Environment directory $(TERRAFORM_DIR) does not exist."; \
    exit 1; \
  fi

# Thêm target tiện ích để check plan nhanh
plan:
  @echo "Planning for environment: $(ENV)..."
  cd $(TERRAFORM_DIR) && terraform plan