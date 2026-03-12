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
      key   = "DATABASE_URL"
      value = var.database_url
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

    env {
      key   = "STRIPE_CHECKOUT_BASE_URL"
      value = var.stripe_checkout_base_url
    }

    env {
      key   = "STRIPE_SECRET_KEY"
      value = var.stripe_secret_key
    }

    env {
      key   = "STRIPE_WEBHOOK_SECRET"
      value = var.stripe_webhook_secret
    }

    env {
      key   = "STRIPE_SUCCESS_URL"
      value = var.stripe_success_url
    }

    env {
      key   = "STRIPE_CANCEL_URL"
      value = var.stripe_cancel_url
    }

    env {
      key   = "STRIPE_MEMBERSHIP_PRICE_ID"
      value = var.stripe_membership_price_id
    }

    env {
      key   = "STRIPE_COURSE_PRICE_ID"
      value = var.stripe_course_price_id
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
