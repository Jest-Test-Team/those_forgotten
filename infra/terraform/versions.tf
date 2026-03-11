terraform {
  required_version = ">= 1.8.0"

  required_providers {
    vercel = {
      source  = "vercel/vercel"
      version = "~> 3.0"
    }

    supabase = {
      source  = "supabase/supabase"
      version = "~> 1.0"
    }

    koyeb = {
      source  = "koyeb/koyeb"
      version = "~> 0.1"
    }
  }
}
