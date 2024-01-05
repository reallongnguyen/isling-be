variable "project" {
  description = "Your GCP Project ID"
  type        = string
}

variable "region" {
  description = "Your project region"
  default     = "us-central1"
  type        = string
}

variable "zone" {
  description = "Your project zone"
  default     = "us-central1-a"
  type        = string
}

variable "env" {
  description = "Your project env"
  default     = "dev"
  type        = string
}

variable "domain" {
  description = "Your domain"
  default     = "dev.isling.me"
  type        = string
}
