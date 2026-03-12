resource "koyeb_app" "platform" {
  name = "${var.project_name}-platform"
}

resource "koyeb_service" "api" {
  app_name = koyeb_app.platform.name

  definition {
    name    = "api-go"
    regions = ["sin"]

    instance_types {
      type = "free"
    }

    scalings {
      min = 1
      max = 1
    }

    ports {
      port     = 8080
      protocol = "http"
    }

    routes {
      path = "/"
      port = 8080
    }

    env {
      key   = "PORT"
      value = "8080"
    }

    env {
      key   = "WEB_ORIGIN"
      value = var.web_origin
    }

    env {
      key   = "ADMIN_EMAILS"
      value = var.admin_emails
    }

    env {
      key   = "INTERNAL_INGEST_TOKEN"
      value = var.internal_ingest_token
    }

    env {
      key   = "SUPABASE_JWT_SECRET"
      value = var.supabase_jwt_secret
    }

    image = var.api_image
  }
}

resource "koyeb_service" "crawler" {
  count    = var.enable_koyeb_crawler_service ? 1 : 0
  app_name = koyeb_app.platform.name

  definition {
    name    = "crawler"
    regions = ["sin"]

    instance_types {
      type = "free"
    }

    scalings {
      min = 0
      max = 1
    }

    env {
      key   = "INGEST_URL"
      value = "${var.api_base_url}/internal/ingest/auctions"
    }

    env {
      key   = "INGEST_TOKEN"
      value = var.internal_ingest_token
    }

    image = var.crawler_image
  }
}
