# Biến mặc định, có thể ghi đè bằng: make deploy ENV=prod
ENV ?= dev
BACKEND_DIR = backend
# Đường dẫn chuẩn từ root đến thư mục môi trường Terraform
TERRAFORM_PATH = infra/envs/$(ENV)

.PHONY: build deploy plan destroy clean

# 1. Build Backend (Gọi Makefile bên trong thư mục backend)
build:
	@echo "--- Building Backend for $(ENV) ---"
	cd $(BACKEND_DIR) && $(MAKE)

clean:
	@echo "--- Cleaning Backend for $(ENV) ---"
	cd $(BACKEND_DIR) && $(MAKE) clean

# 2. Deploy: Tự động chạy build -> cd vào đúng môi trường -> apply
deploy: build
	@echo "--- Deploying Infrastructure to [$(ENV)] ---"
	@if [ -d "$(TERRAFORM_PATH)" ]; then \
		cd $(TERRAFORM_PATH) && terraform init && terraform apply -auto-approve; \
	else \
		echo "Error: Directory $(TERRAFORM_PATH) not found!"; \
		exit 1; \
	fi

# 3. Plan: Chỉ xem thay đổi mà không apply
plan:
	@echo "--- Planning Infrastructure for [$(ENV)] ---"
	cd $(TERRAFORM_PATH) && terraform plan

destroy:
	@echo "--- Destroying Infrastructure to [$(ENV)] ---"
	@if [ -d "$(TERRAFORM_PATH)" ]; then \
		cd $(TERRAFORM_PATH) && terraform destroy -auto-approve; \
	else \
		echo "Error: Directory $(TERRAFORM_PATH) not found!"; \
		exit 1; \
	fi