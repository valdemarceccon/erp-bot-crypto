variable "VERSION" {
  default = "latest"
}

group "default" {
  targets = ["proxy", "frontend", "backend", "daemon"]
}

target "proxy" {
  context = "nginx-conf"
  dockerfile = "Dockerfile"
  tags = ["ghcr.io/valdemarceccon/proxy:${VERSION}"]
}

target "frontend" {
  context = "frontend"
  dockerfile = "Dockerfile.prod"
  tags = ["ghcr.io/valdemarceccon/frontend:${VERSION}"]
}

target "backend" {
  context = "backend"
  dockerfile = "Dockerfile.api"
  tags = ["ghcr.io/valdemarceccon/backend:${VERSION}"]
}

target "daemon" {
  context = "backend"
  dockerfile = "Dockerfile.daemon"
  tags = ["ghcr.io/valdemarceccon/daemon:${VERSION}"]
}
