resource "vercel_project" "web" {
  name           = "${var.project_name}-web"
  framework      = "nextjs"
  root_directory = "apps/web"

  git_repository = {
    type = "github"
    repo = var.github_repo
  }
}

resource "vercel_project_environment_variable" "next_public_api_base_url" {
  project_id = vercel_project.web.id
  key        = "NEXT_PUBLIC_API_BASE_URL"
  value      = var.api_base_url
  target     = ["production", "preview", "development"]
}

resource "vercel_project_environment_variable" "next_public_site_url" {
  project_id = vercel_project.web.id
  key        = "NEXT_PUBLIC_SITE_URL"
  value      = var.site_url
  target     = ["production", "preview", "development"]
}

resource "vercel_project_domain" "web_domain" {
  project_id = vercel_project.web.id
  domain     = var.web_domain
}
