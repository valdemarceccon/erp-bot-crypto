variable "VERSION" {
  default = "latest"
}

group "default" {
  targets = ["proxy", "frontend", "backend", "daemon"]
}

target "proxy" {
  context = "nginx-conf"
  dockerfile = "Dockerfile"
  tags = ["ghcr.io/valdemarceccon/proxy:${VERSION}", "ghcr.io/valdemarceccon/frontend:latest"]
}

target "frontend" {
  context = "frontend"
  dockerfile = "Dockerfile.prod"
  tags = ["ghcr.io/valdemarceccon/frontend:${VERSION}", "ghcr.io/valdemarceccon/frontend:latest"]
}

target "backend" {
  context = "backend"
  dockerfile = "docker/Dockerfile.api.prod"
  tags = ["ghcr.io/valdemarceccon/backend:${VERSION}", "ghcr.io/valdemarceccon/frontend:latest"]
}

target "daemon" {
  context = "backend"
  dockerfile = "docker/Dockerfile.daemon"
  tags = ["ghcr.io/valdemarceccon/daemon:${VERSION}", "ghcr.io/valdemarceccon/frontend:latest"]
}
