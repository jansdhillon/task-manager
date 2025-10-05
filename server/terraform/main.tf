provider "google" {
  project = var.gcp_project_id
  region  = var.region
}

resource "google_cloud_run_v2_service" "task_manager" {
    name = var.service_name
    location = var.region
    deletion_protection = false
    ingress = "INGRESS_TRAFFIC_ALL"

    template {
        containers {
            image = "gcr.io/${var.artifact_repo}"
        }
    }
}
