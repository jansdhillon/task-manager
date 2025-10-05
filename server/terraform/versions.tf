terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "7.5.0"
    }
  }
  required_version = ">= 1.10"
  backend "gcs" {
    bucket                      = "task-manager-terraform-state"
    prefix                      = "prod"
    impersonate_service_account = "941741555207-compute@developer.gserviceaccount.com"
  }
}
