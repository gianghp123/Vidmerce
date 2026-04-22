.ONESHELL:
SHELL := /bin/bash
# --- Configuration ---
ENV             ?= dev
BACKEND_DIR     = backend
FRONTEND_DIR    = frontend
SHARED_DIR      = shared
SCHEMAS_DIR     = $(SHARED_DIR)/schemas
SCRIPTS_DIR     = $(SHARED_DIR)/scripts

TERRAFORM_PATH  = infra/envs/$(ENV)

# --- Go Type Generation Config ---
GO_ENUMS_OUT    = $(BACKEND_DIR)/internal/core/enums
GO_MODELS_OUT   = $(BACKEND_DIR)/internal/database/models
GO_ENUMS_IMPORT = github.com/gianghp123/Vidmerce/backend/internal/core/enums

# --- TS Type Generation Config ---
TS_ENUMS_OUT    = $(FRONTEND_DIR)/src/lib/enums
TS_MODELS_OUT   = $(FRONTEND_DIR)/src/lib/models

# --- Terraform Settings ---
TF_CMD          = terraform
TF_FLAGS        ?=
TF_PARALLELISM  ?= 5
TF_DEBUG        ?= 0

ifeq ($(TF_DEBUG),1)
export TF_LOG      = DEBUG
export TF_LOG_PATH = terraform.log
endif

.PHONY: build deploy plan destroy generate-types

# --- Type Generation ---
generate-types:
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

# --- Main Logic ---
build:
	@echo "--- Building Backend for $(ENV) ---"
	cd $(BACKEND_DIR) && $(MAKE)

	@echo "--- Building Frontend for $(ENV) ---"
	cd ../$(FRONTEND_DIR) && npm run build

plan:
	@cd $(TERRAFORM_PATH) && $(TF_CMD) plan -parallelism=$(TF_PARALLELISM) $(TF_FLAGS)

deploy:
	@echo "--- Deploying to AWS [$(ENV)] ---"
	@cd $(TERRAFORM_PATH) && $(TF_CMD) init && $(TF_CMD) apply -auto-approve -parallelism=$(TF_PARALLELISM) $(TF_FLAGS)

destroy:
	@echo "--- Destroying AWS [$(ENV)] ---"
	@cd $(TERRAFORM_PATH) && $(TF_CMD) init && $(TF_CMD) destroy -auto-approve -parallelism=$(TF_PARALLELISM) $(TF_FLAGS)