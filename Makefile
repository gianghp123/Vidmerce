.ONESHELL:
SHELL := /bin/bash
# --- Configuration ---
ENV             ?= dev
BACKEND_DIR     = backend
TERRAFORM_PATH  = infra/envs/$(ENV)
LOCAL_OVERRIDE  = $(TERRAFORM_PATH)/override.tf.json

# --- LocalStack & Endpoints ---
# Thay đổi URL ở đây nếu bạn chạy LocalStack trên máy khác hoặc Docker Network
LOCALSTACK_HOST = localhost
LOCALSTACK_PORT = 4566
LOCALSTACK_URL  = http://$(LOCALSTACK_HOST):$(LOCALSTACK_PORT)
S3_ENDPOINT     = http://s3.localhost.localstack.cloud:$(LOCALSTACK_PORT)
# S3_ENDPOINT     = http://$(LOCALSTACK_HOST):$(LOCALSTACK_PORT)
# --- Terraform Settings ---
TF_CMD          = terraform
TF_FLAGS        ?=
TF_PARALLELISM  ?= 5
TF_DEBUG        ?= 0

ifeq ($(TF_DEBUG),1)
export TF_LOG      = DEBUG
export TF_LOG_PATH = terraform.log
endif

.PHONY: build deploy plan destroy clean setup-local setup-cloud test-lambda-local

# --- Helpers ---

# Sử dụng printf để định dạng JSON sạch sẽ và dễ đọc hơn
setup-local:
	@PERSISTENCE=1 localstack start -e SERVICES=s3,dynamodb,lambda,iam -d
	@cat > $(LOCAL_OVERRIDE) <<'EOF'
	{
	  "variable": {
			"is_local": { "default": true }
	  },
	  "provider": {
	    "aws": {
	      "access_key": "test",
	      "secret_key": "test",
	      "skip_credentials_validation": true,
	      "skip_metadata_api_check": true,
	      "skip_requesting_account_id": true,
	      "endpoints": {
	        "dynamodb": "$(LOCALSTACK_URL)",
	        "lambda":   "$(LOCALSTACK_URL)",
	        "s3":       "$(S3_ENDPOINT)",
	        "iam":      "$(LOCALSTACK_URL)",
	        "sts":      "$(LOCALSTACK_URL)"
# 					"apigatewayv2": "$(LOCALSTACK_URL)"
	      }
	    }
	  }
	}
	EOF
	@echo "--- [LOCAL] LocalStack override created at $(LOCAL_OVERRIDE) ---"

setup-cloud:
	@localstack stop
	@rm -f $(LOCAL_OVERRIDE)
	@echo "--- [CLOUD] LocalStack override removed. Ready for Real AWS. ---"

# --- Main Logic ---

build:
	@echo "--- Building Backend for $(ENV) ---"
	cd $(BACKEND_DIR) && $(MAKE)

clean:
	cd $(BACKEND_DIR) && $(MAKE) clean

plan:
	@cd $(TERRAFORM_PATH) && $(TF_CMD) plan -parallelism=$(TF_PARALLELISM) $(TF_FLAGS)

deploy: setup-cloud build
	@echo "--- Deploying to [$(ENV)] ---"
	cd $(TERRAFORM_PATH) && $(TF_CMD) init && $(TF_CMD) apply -auto-approve -parallelism=$(TF_PARALLELISM) $(TF_FLAGS)

deploy-local: setup-local
	@echo "--- Deploying to LocalStack ([$(ENV)] config) ---"
	cd $(TERRAFORM_PATH) && $(TF_CMD) init && $(TF_CMD) apply -auto-approve -parallelism=$(TF_PARALLELISM) $(TF_FLAGS)

destroy: 
	cd $(TERRAFORM_PATH) && $(TF_CMD) destroy -auto-approve -parallelism=$(TF_PARALLELISM) $(TF_FLAGS)


destroy-local:
	@$(MAKE) setup-local
	@cd $(TERRAFORM_PATH) && $(TF_CMD) destroy -auto-approve -parallelism=$(TF_PARALLELISM) $(TF_FLAGS)
	@$(MAKE) setup-cloud