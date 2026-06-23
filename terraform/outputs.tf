output "cluster_endpoint" {
  description = "EKS cluster endpoint"
  value       = module.eks.cluster_endpoint
}

output "cluster_name" {
  description = "EKS cluster name"
  value       = module.eks.cluster_name
}

output "rds_endpoint" {
  description = "RDS endpoint"
  value       = module.rds.db_instance_endpoint
  sensitive   = true
}

output "redis_endpoint" {
  description = "ElastiCache Redis endpoint"
  value       = module.elasticache.primary_endpoint_address
}

output "msk_bootstrap_brokers" {
  description = "MSK bootstrap brokers"
  value       = module.msk.bootstrap_brokers_tls
  sensitive   = true
}

output "s3_documents_bucket" {
  description = "S3 bucket for documents"
  value       = module.s3_documents.s3_bucket_id
}

output "db_password_secret_arn" {
  description = "ARN of the database password in Secrets Manager"
  value       = aws_secretsmanager_secret.db_password.arn
}
