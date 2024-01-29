provider "google" {
  project = var.project
  region  = var.region
}

data "google_client_config" "current" {
  provider = google
}

resource "google_cloud_run_v2_service" "role-finder" {
  name         = "role-finder"
  location     = "europe-west4"
  ingress      = "INGRESS_TRAFFIC_ALL"
  provider     = google-beta
  launch_stage = "BETA"

  template {
    execution_environment = "EXECUTION_ENVIRONMENT_GEN2"
    containers {
      image = "eu.gcr.io/binx-io-public/gcp-role-finder:latest"
      command = ["/app/gcp-role-finder", "serve", "--from-file", "--data-file", "/mnt/data/roles.json"]
      volume_mounts {
        name       = "bucket"
        mount_path = "/mnt/data"
      }
    }

    volumes {
      name = "bucket"
      gcs {
        bucket = google_storage_bucket.role-cache.name
      }
    }
  }
}

resource "google_cloud_run_service_iam_binding" "role-finder-invoker" {
  location = google_cloud_run_v2_service.role-finder.location
  project  = google_cloud_run_v2_service.role-finder.project
  service  = google_cloud_run_v2_service.role-finder.name
  role     = "roles/run.invoker"
  members  = ["allUsers"]
}

resource "google_service_account" "role-finder" {
  account_id = "role-finder"
}

resource "google_project_iam_member" "role-finder" {
  for_each = toset(["roles/iam.roleViewer"])
  role     = each.value
  member   = google_service_account.role-finder.member
  project  = google_service_account.role-finder.project
}

resource "google_storage_bucket" "role-cache" {
  name     = format("role-finder-cache-%s-%s", var.region, data.google_client_config.current.project)
  location = var.region

  uniform_bucket_level_access = true
}

resource "google_storage_bucket_object" "roles" {
  bucket       = google_storage_bucket.role-cache.name
  name         = "roles.json"
  content_type = "application/json"
  content      = file(format("%s/../data/roles.json", path.module))
}

resource "google_storage_bucket_iam_binding" "owners" {
  bucket = google_storage_bucket.role-cache.name
  role   = "roles/storage.objectUser"
  members = [
    google_service_account.role-finder.member
  ]
}
