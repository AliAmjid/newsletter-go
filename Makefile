.PHONY: help terraform-plan terraform-apply

ENV_FILE := .env

# Export .env into environment variables
include $(ENV_FILE)
export

# Define a variable for the folder route
folder ?= .

help:
	@echo "Available commands:"
	@echo "  make terraform-plan        - Run Terraform plan command"
	@echo "  make terraform-apply       - Run Terraform apply command"


terraform-plan:
	dotenvx run -fk .env.keys -f $(folder)/.env -- terraform -chdir=$(folder) plan

terraform-apply:
	dotenvx run -fk .env.keys -f $(folder)/.env -- terraform -chdir=$(folder) apply
