
provider "google" {
  version = "~> 3.9.0"
}


# Zone Info for roger-beathard.com
resource "google_dns_managed_zone" "roger-beathard-com" {
  name          = "roger-beathard-com"
  dns_name      = "roger-beathard.com."
  description   = "zone file for roger-beathard.com"
  project       = "my-dns-279215"
}

# Zone Record for www.roger-beathard.com
resource "google_dns_record_set" "www-roger-beathard-com" {
  managed_zone  = "roger-beathard-com"
  project       = "my-dns-279215"

  name          = "www.roger-beathard.com."
  type          = "A"
  rrdatas       = ["10.1.2.1"]
  ttl           = "300"
}


# Zone Record for azure.roger-beathard.com
resource "google_dns_record_set" "azure-roger-beathard-com" {
  managed_zone  = "roger-beathard-com"
  project       = "my-dns-279215"

  name          = "azure.roger-beathard.com."
  type          = "TXT"
  rrdatas       = ["\"v=spf1 ip4:111.111.111.111 include:backoff.example.com -all\""]
  ttl           = "300"
}


# Zone Record for home-red.roger-beathard.com
resource "google_dns_record_set" "home-red-roger-beathard-com" {
  managed_zone  = "roger-beathard-com"
  project       = "my-dns-279215"

  name          = "home-red.roger-beathard.com."
  type          = "A"
  rrdatas       = ["10.1.4.3"]
  ttl           = "300"
}


# Zone Record for home-green.roger-beathard.com
resource "google_dns_record_set" "home-green-roger-beathard-com" {
  managed_zone  = "roger-beathard-com"
  project       = "my-dns-279215"

  name          = "home-green.roger-beathard.com."
  type          = "A"
  rrdatas       = ["10.1.4.3"]
  ttl           = "300"
}


# Zone Record for aweb.roger-beathard.com
resource "google_dns_record_set" "aweb-roger-beathard-com" {
  managed_zone  = "roger-beathard-com"
  project       = "my-dns-279215"

  name          = "aweb.roger-beathard.com."
  type          = "A"
  rrdatas       = ["10.1.4.3"]
  ttl           = "300"
}


# Zone Record for apple.roger-beathard.com
resource "google_dns_record_set" "apple-roger-beathard-com" {
  managed_zone  = "roger-beathard-com"
  project       = "my-dns-279215"

  name          = "apple.roger-beathard.com."
  type          = "A"
  rrdatas       = ["10.1.5.6"]
  ttl           = "300"
}


# Zone Record for home-lb.roger-beathard.com
resource "google_dns_record_set" "home-lb-roger-beathard-com" {
  managed_zone  = "roger-beathard-com"
  project       = "my-dns-279215"

  name          = "home-lb.roger-beathard.com."
  type          = "CNAME"
  rrdatas       = ["home-red.rvbtech.com."]
  ttl           = "300"
}


# Zone Record for home2.roger-beathard.com
resource "google_dns_record_set" "home2-roger-beathard-com" {
  managed_zone  = "roger-beathard-com"
  project       = "my-dns-279215"

  name          = "home2.roger-beathard.com."
  type          = "CNAME"
  rrdatas       = ["home-red.vbtech.com."]
  ttl           = "300"
}

