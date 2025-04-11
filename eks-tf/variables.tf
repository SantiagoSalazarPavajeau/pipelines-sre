variable "cluster_name" {
  description = "Name of the mock cluster"
  type        = string
  default     = "test-mock-cluster"
}

variable "region" {
  description = "Mock AWS region"
  type        = string
  default     = "us-mock-1"
}

variable "node_group_name" {
  description = "Mock node group name"
  type        = string
  default     = "test-node-group"
}
