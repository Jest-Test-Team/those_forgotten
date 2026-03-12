variable "project_name" {
  type        = string
  default     = "those-forgotten"
  description = "Base project name used across Vercel, Supabase, and Koyeb."
}

variable "github_repo" {
  type        = string
  description = "GitHub repository in owner/repo format."
}

variable "vercel_api_token" {
  type        = string
  sensitive   = true
  description = "Vercel API token."
}

variable "vercel_team_id" {
  type        = string
  default     = null
  description = "Optional Vercel team id."
}

variable "web_domain" {
  type        = string
  default     = "customs-auction.example.com"
  description = "Primary public domain for the web app."
}

variable "api_base_url" {
  type        = string
  default     = "https://api-customs-auction.example.com"
  description = "Public API base URL used by the Next.js app."
}

variable "site_url" {
  type        = string
  default     = "https://customs-auction.example.com"
  description = "Primary site URL used by OAuth and metadata."
}

variable "supabase_access_token" {
  type        = string
  sensitive   = true
  description = "Supabase access token."
}

variable "supabase_organization_id" {
  type        = string
  description = "Supabase organization id."
}

variable "supabase_region" {
  type        = string
  default     = "ap-northeast-1"
  description = "Supabase project region."
}

variable "db_pass" {
  type        = string
  sensitive   = true
  description = "Supabase database password."
}

variable "koyeb_token" {
  type        = string
  sensitive   = true
  description = "Koyeb API token."
}

variable "koyeb_org_id" {
  type        = string
  description = "Koyeb organization id."
}

variable "api_image" {
  type        = string
  default     = "ghcr.io/jest-test-team/those_forgotten-api:latest"
  description = "Container image for the Go API."
}

variable "crawler_image" {
  type        = string
  default     = "ghcr.io/jest-test-team/those_forgotten-crawler:latest"
  description = "Container image for the crawler."
}

variable "enable_koyeb_crawler_service" {
  type        = bool
  default     = false
  description = "Enable crawler deployment on Koyeb. Keep false on strict free-tier setups and use GitHub Actions cron instead."
}

variable "internal_ingest_token" {
  type        = string
  sensitive   = true
  description = "Shared secret for crawler -> API ingest."
}

variable "supabase_jwt_secret" {
  type        = string
  sensitive   = true
  description = "JWT secret used by api-go to verify Supabase access tokens."
}

variable "web_origin" {
  type        = string
  default     = "https://customs-auction.example.com"
  description = "Allowed web origin for API CORS."
}

variable "admin_emails" {
  type        = string
  default     = "admin@example.com"
  description = "Comma-separated admin allowlist consumed by the web SSR role guard."
}
