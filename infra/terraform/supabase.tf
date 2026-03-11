resource "supabase_project" "main" {
  organization_id   = var.supabase_organization_id
  name              = "${var.project_name}-supabase"
  region            = var.supabase_region
  database_password = var.db_pass
}

resource "supabase_project_api_settings" "main" {
  project_ref = supabase_project.main.id

  db_schema            = "public,storage,graphql_public"
  max_rows             = 1000
  jwt_expiry           = 3600
  enable_refresh_token = true
}

resource "supabase_project_auth_settings" "main" {
  project_ref           = supabase_project.main.id
  site_url              = var.site_url
  additional_redirect_urls = [
    "${var.site_url}/auth/callback",
    "http://localhost:3000/auth/callback",
  ]
  jwt_expiry = 3600
}
