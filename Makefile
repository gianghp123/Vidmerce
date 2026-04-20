.ONESHELL:
SHELL := /bin/bash
# --- Configuration ---
ENV             ?= dev
BACKEND_DIR     = backend
FRONTEND_DIR		= frontend
SHARED_DIR      = shared
SCHEMAS_DIR     = $(SHARED_DIR)/schemas
SCRIPTS_DIR     = $(SHARED_DIR)/scripts


TERRAFORM_PATH  = infra/envs/$(ENV)
LOCAL_OVERRIDE  = $(TERRAFORM_PATH)/override.tf.json

# --- Go Type Generation Config ---
GO_ENUMS_OUT    = $(BACKEND_DIR)/internal/core/enums
GO_MODELS_OUT   = $(BACKEND_DIR)/internal/database/models
GO_ENUMS_IMPORT = github.com/gianghp123/Vidmerce/backend/internal/core/enums

# --- TS Type Generation Config ---
TS_ENUMS_OUT    = $(FRONTEND_DIR)/src/lib/enums
TS_MODELS_OUT   = $(FRONTEND_DIR)/src/lib/models


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

.PHONY: build deploy plan destroy setup-local setup-cloud test-lambda-local generate

# --- Type Generation ---
generate:
	@echo "--- [GENERATE] Generating Go Models & Enums ---"
	@python3 $(SCRIPTS_DIR)/generate_go_types.py \
		--schema-dir $(SCHEMAS_DIR) \
		--enums-dir $(GO_ENUMS_OUT) \
		--models-dir $(GO_MODELS_OUT) \
		--enums-import $(GO_ENUMS_IMPORT)

	@echo "--- [GENERATE] Generating TS Models & Enums ---"
	@python3 $(SCRIPTS_DIR)/generate_ts_types.py \
		--schema-dir $(SCHEMAS_DIR) \
		--enums-dir $(TS_ENUMS_OUT) \
		--models-dir $(TS_MODELS_OUT)
	@echo "--- [GENERATE] All types updated successfully ---"

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

	@echo "--- Building Frontend for $(ENV) ---"
	cd ../$(FRONTEND_DIR) && npm run build

plan:
	@cd $(TERRAFORM_PATH) && $(TF_CMD) plan -parallelism=$(TF_PARALLELISM) $(TF_FLAGS)

deploy: setup-cloud build
	@echo "--- Deploying to [$(ENV)] ---"
	cd $(TERRAFORM_PATH) && $(TF_CMD) init && $(TF_CMD) apply -auto-approve -parallelism=$(TF_PARALLELISM) $(TF_FLAGS)

deploy-local: setup-local build
	@echo "--- Deploying to LocalStack ([$(ENV)] config) ---"
	cd $(TERRAFORM_PATH) && $(TF_CMD) init && $(TF_CMD) apply -auto-approve -parallelism=$(TF_PARALLELISM) $(TF_FLAGS)

destroy: 
	cd $(TERRAFORM_PATH) && $(TF_CMD) destroy -auto-approve -parallelism=$(TF_PARALLELISM) $(TF_FLAGS)


destroy-local:
	@$(MAKE) setup-local
	@cd $(TERRAFORM_PATH) && $(TF_CMD) destroy -auto-approve -parallelism=$(TF_PARALLELISM) $(TF_FLAGS)
	@$(MAKE) setup-cloud