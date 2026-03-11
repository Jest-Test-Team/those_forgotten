output "vercel_project_id" {
  value = vercel_project.web.id
}

output "supabase_project_id" {
  value = supabase_project.main.id
}

output "koyeb_app_name" {
  value = koyeb_app.platform.name
}
