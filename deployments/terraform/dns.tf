data "google_dns_managed_zone" "dns_zone" {
  project = var.google_project
  name    = var.dns_zone_name
}

resource "google_dns_record_set" "home_dns_entry" {
  name    = "home.${data.google_dns_managed_zone.dns_zone.dns_name}"
  type    = "CNAME"
  ttl     = 300
  project = var.google_project

  managed_zone = data.google_dns_managed_zone.dns_zone.name

  rrdatas = ["${var.dyn_dns_name}."]
}

resource "google_dns_record_set" "grafana_dns_entry" {
  name    = "grafana.${data.google_dns_managed_zone.dns_zone.dns_name}"
  type    = "CNAME"
  ttl     = 300
  project = var.google_project

  managed_zone = data.google_dns_managed_zone.dns_zone.name

  rrdatas = [google_dns_record_set.home_dns_entry.name]
}

resource "google_dns_record_set" "experimental_dns_entry" {
  name    = "experimental.${data.google_dns_managed_zone.dns_zone.dns_name}"
  type    = "CNAME"
  ttl     = 300
  project = var.google_project

  managed_zone = data.google_dns_managed_zone.dns_zone.name

  rrdatas = [google_dns_record_set.home_dns_entry.name]
}

