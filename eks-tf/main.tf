resource "null_resource" "mock_eks_cluster" {
  triggers = {
    cluster_name = var.cluster_name
    region       = var.region
  }

  provisioner "local-exec" {
    command = "echo 'Pretending to provision EKS cluster named ${var.cluster_name}'"
  }
}

resource "null_resource" "mock_node_group" {
  triggers = {
    node_group = var.node_group_name
  }

  provisioner "local-exec" {
    command = "echo 'Mock node group: ${var.node_group_name}'"
  }
}

terraform {
  required_providers {
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = "~> 2.0"
    }
    helm = {
      source = "hashicorp/helm"
      version = "~> 2.4"
    }
  }

  required_version = ">= 1.3.0"
}

provider "kubernetes" {
  config_path = "~/.kube/config"
}

provider "helm" {
  kubernetes {
    config_path = "~/.kube/config"
  }
}

resource "helm_release" "services" {
  for_each = {
    redis-store = {
      image_repo = "redis"
      image_tag  = "latest"
      port       = 6379
    }
    backend-app = {
      image_repo = "kennethreitz/httpbin"
      image_tag  = "latest"
      port       = 80
    }
    stream-service = {
      image_repo = "nats"
      image_tag  = "latest"
      port       = 4222
    }
    api-gateway = {
      image_repo = "nginx"
      image_tag  = "latest"
      port       = 80
    }
  }

  name       = each.key
  chart      = "../helm"
  namespace  = "default"

  set {
    name  = "image.repository"
    value = each.value.image_repo
  }

  set {
    name  = "image.tag"
    value = each.value.image_tag
  }

  set {
    name  = "service.port"
    value = each.value.port
  }

  set {
    name  = "deployment.name"
    value = each.key
  }
}
