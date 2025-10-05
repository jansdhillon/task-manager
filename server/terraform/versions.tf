terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "7.5.0"
    }
  }
  required_version = ">= 1.10"
  backend "gcs" {}
}
