variable "google_project" {
  type        = string
  description = "The GCP project to create resources in."
}

variable "dns_zone_name" {
  type        = string
  description = "The name of the DNS zone to create the record in."
}

variable "dyn_dns_name" {
  type        = string
  description = "The dynamic DNS name to point to."
}
