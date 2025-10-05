variable "gcp_project_id" {
  type = string
}

variable "region" {
  type = string
}

variable "artifact_repo" {
  type = string
}

variable "service_name" {
  type = string
}

variable "allow_unauthenticated" {
  type    = bool
  default = true
}
