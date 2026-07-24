terraform {
  backend "s3" {
    # Backend configuration provided via backend config file (backend.tfvars)
    # terraform init -backend-config=backend.tfvars
  }
}