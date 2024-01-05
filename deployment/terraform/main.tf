terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "4.51.0"
    }
  }
}

provider "google" {
  // credentials = file("<NAME>.json")

  project = var.project
  region  = var.region
  zone    = var.zone
}

resource "google_compute_network" "vpc_network" {
  name = "iplay-${var.env}-vpc"
}

resource "google_compute_firewall" "iplay_allow_ssh" {
  name    = "iplay-${var.env}-vpc-allow-ssh"
  network = google_compute_network.vpc_network.name

  allow {
    protocol = "tcp"
    ports    = ["22"]
  }

  source_ranges = ["0.0.0.0/0"]
}

resource "google_compute_firewall" "iplay_allow_http" {
  name    = "iplay-${var.env}-vpc-allow-http"
  network = google_compute_network.vpc_network.name

  allow {
    protocol = "tcp"
    ports    = ["80"]
  }

  source_ranges = ["0.0.0.0/0"]
  target_tags   = ["http-server"]
}

resource "google_compute_firewall" "iplay_allow_https" {
  name    = "iplay-${var.env}-vpc-allow-https"
  network = google_compute_network.vpc_network.name

  allow {
    protocol = "tcp"
    ports    = ["443"]
  }

  source_ranges = ["0.0.0.0/0"]
  target_tags   = ["https-server"]
}

resource "google_compute_firewall" "iplay_allow_internal" {
  name          = "iplay-${var.env}-vpc-allow-internal"
  network       = google_compute_network.vpc_network.name
  priority      = 65534
  direction     = "INGRESS"
  source_ranges = ["10.128.0.0/9"]

  allow {
    protocol = "tcp"
    ports    = ["0-65535"]
  }

  allow {
    protocol = "udp"
    ports    = ["0-65535"]
  }

  allow {
    protocol = "icmp"
  }
}

# NAT
resource "google_compute_subnetwork" "subnet01" {
  name          = "iplay-${var.env}-subnet01"
  network       = google_compute_network.vpc_network.id
  ip_cidr_range = "10.0.0.0/16"
  region        = var.region
}

resource "google_compute_router" "router01" {
  name    = "iplay-${var.env}-router01"
  region  = google_compute_subnetwork.subnet01.region
  network = google_compute_network.vpc_network.id

  bgp {
    asn = 64514
  }
}

resource "google_compute_router_nat" "nat01" {
  name                               = "iplay-${var.env}-router-nat"
  router                             = google_compute_router.router01.name
  region                             = google_compute_router.router01.region
  nat_ip_allocate_option             = "AUTO_ONLY"
  source_subnetwork_ip_ranges_to_nat = "ALL_SUBNETWORKS_ALL_IP_RANGES"

  log_config {
    enable = true
    filter = "ERRORS_ONLY"
  }
}

# API01
resource "google_compute_address" "vm_instance_api01_address_int" {
  name         = "iplay-${var.env}-api01-ip-int"
  address_type = "INTERNAL"
  address      = "10.128.0.7"
}

resource "google_compute_instance" "vm_instance_api01" {
  name         = "iplay-${var.env}-api01"
  machine_type = "e2-micro"
  tags         = ["http-server", "https-server", "lb-health-check"]

  boot_disk {
    initialize_params {
      image = "debian-cloud/debian-12"
    }
  }

  network_interface {
    network    = google_compute_network.vpc_network.name
    network_ip = google_compute_address.vm_instance_api01_address_int.address
  }

  # Because provisioning model is SPOT, the VM may be terminated at any time
  # You should run the /deployment/scripts/wakeup-vm.sh every 1 minute
  # Use crontab to schedule this task
  scheduling {
    automatic_restart           = false
    on_host_maintenance         = "TERMINATE"
    preemptible                 = true
    provisioning_model          = "SPOT"
    instance_termination_action = "STOP"
  }
}

# Gorse
resource "google_compute_address" "vm_instance_gorse_address" {
  name         = "iplay-${var.env}-gorse-ip"
  address_type = "INTERNAL"
  address      = "10.128.0.6"
}

resource "google_compute_instance" "vm_instance_gorse" {
  name         = "iplay-${var.env}-gorse"
  machine_type = "e2-micro"

  boot_disk {
    initialize_params {
      image = "debian-cloud/debian-12"
    }
  }

  network_interface {
    network    = google_compute_network.vpc_network.name
    network_ip = google_compute_address.vm_instance_gorse_address.address
  }
}

# Surreal
resource "google_compute_address" "vm_instance_surreal_address_int" {
  name         = "iplay-${var.env}-surreal-ip-int"
  address_type = "INTERNAL"
  address      = "10.128.0.8"
}

resource "google_compute_instance" "vm_instance_surreal" {
  name         = "iplay-${var.env}-surreal"
  machine_type = "e2-micro"
  tags         = ["http-server", "https-server"]

  boot_disk {
    initialize_params {
      image = "debian-cloud/debian-12"
    }
  }

  network_interface {
    network    = google_compute_network.vpc_network.name
    network_ip = google_compute_address.vm_instance_surreal_address_int.address
  }
}

# ELB
resource "google_compute_address" "vm_instance_elb_address_ext" {
  name = "iplay-${var.env}-elb-ip"
}

resource "google_compute_instance" "vm_instance_elb" {
  name         = "iplay-${var.env}-elb"
  machine_type = "e2-micro"
  tags         = ["http-server", "https-server"]

  boot_disk {
    initialize_params {
      image = "debian-cloud/debian-12"
    }
  }

  network_interface {
    network = google_compute_network.vpc_network.name
    access_config {
      nat_ip = google_compute_address.vm_instance_elb_address_ext.address
    }
  }
}

# DNS
resource "google_dns_managed_zone" "dns_zone" {
  name        = "iplay-${var.env}-dns"
  dns_name    = "${var.domain}."
  description = "iplay-${var.env} DNS zone"
}

resource "google_dns_record_set" "dns_record_elb" {
  name         = "*.${google_dns_managed_zone.dns_zone.dns_name}"
  managed_zone = google_dns_managed_zone.dns_zone.name
  type         = "A"
  ttl          = 300

  rrdatas = [google_compute_address.vm_instance_elb_address_ext.address]
}
