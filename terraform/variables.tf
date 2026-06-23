variable "environment" {
  description = "Deployment environment"
  type        = string
  default     = "prod"
}

variable "aws_region" {
  description = "AWS region"
  type        = string
  default     = "us-east-1"
}

variable "vpc_cidr" {
  description = "VPC CIDR block"
  type        = string
  default     = "10.0.0.0/16"
}

variable "cluster_version" {
  description = "Kubernetes version"
  type        = string
  default     = "1.28"
}

variable "domain_name" {
  description = "Domain name for the platform"
  type        = string
  default     = "platform.example.com"
}

variable "ssl_cert_arn" {
  description = "ACM SSL certificate ARN"
  type        = string
  sensitive   = true
}
