terraform {
  required_version = ">= 1.6"
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.30"
    }
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = "~> 2.24"
    }
    helm = {
      source  = "hashicorp/helm"
      version = "~> 2.12"
    }
    random = {
      source  = "hashicorp/random"
      version = "~> 3.5"
    }
  }
  backend "s3" {
    bucket         = "platform-terraform-state"
    key            = "platform/terraform.tfstate"
    region         = "us-east-1"
    encrypt        = true
    dynamodb_table = "platform-terraform-locks"
  }
}

provider "aws" {
  region = var.aws_region

  default_tags {
    tags = {
      Environment = var.environment
      Project     = "document-intelligence-platform"
      ManagedBy   = "terraform"
    }
  }
}

provider "kubernetes" {
  host                   = module.eks.cluster_endpoint
  cluster_ca_certificate = base64decode(module.eks.cluster_certificate_authority_data)
  exec {
    api_version = "client.authentication.k8s.io/v1beta1"
    command     = "aws"
    args        = ["eks", "get-token", "--cluster-name", module.eks.cluster_name]
  }
}

provider "helm" {
  kubernetes {
    host                   = module.eks.cluster_endpoint
    cluster_ca_certificate = base64decode(module.eks.cluster_certificate_authority_data)
    exec {
      api_version = "client.authentication.k8s.io/v1beta1"
      command     = "aws"
      args        = ["eks", "get-token", "--cluster-name", module.eks.cluster_name]
    }
  }
}

data "aws_availability_zones" "available" {
  state = "available"
}

locals {
  name_prefix = "platform-${var.environment}"
  azs         = slice(data.aws_availability_zones.available.names, 0, 3)
}

module "vpc" {
  source  = "terraform-aws-modules/vpc/aws"
  version = "5.5.1"

  name = "${local.name_prefix}-vpc"
  cidr = "10.0.0.0/16"

  azs             = local.azs
  private_subnets = ["10.0.1.0/24", "10.0.2.0/24", "10.0.3.0/24"]
  public_subnets  = ["10.0.101.0/24", "10.0.102.0/24", "10.0.103.0/24"]

  enable_nat_gateway   = true
  enable_vpn_gateway   = false
  enable_dns_hostnames = true

  public_subnet_tags = {
    "kubernetes.io/role/elb" = "1"
  }
  private_subnet_tags = {
    "kubernetes.io/role/internal-elb" = "1"
  }
}

module "eks" {
  source  = "terraform-aws-modules/eks/aws"
  version = "19.20.0"

  cluster_name    = "${local.name_prefix}-cluster"
  cluster_version = "1.28"

  vpc_id     = module.vpc.vpc_id
  subnet_ids = module.vpc.private_subnets

  cluster_endpoint_public_access = true

  eks_managed_node_groups = {
    platform = {
      desired_size = 3
      min_size     = 3
      max_size     = 10

      instance_types = ["m6i.large", "m6a.large"]
      capacity_type  = "ON_DEMAND"

      block_device_mappings = {
        xvda = {
          device_name = "/dev/xvda"
          ebs = {
            volume_size           = 100
            volume_type           = "gp3"
            iops                  = 3000
            throughput            = 125
            delete_on_termination = true
          }
        }
      }
    }

    ai = {
      desired_size = 1
      min_size     = 1
      max_size     = 5

      instance_types = ["g5.xlarge", "p3.2xlarge"]
      capacity_type  = "SPOT"

      labels = {
        workload = "ai-inference"
      }
      taints = {
        gpu = {
          key    = "nvidia.com/gpu"
          value  = "true"
          effect = "NO_SCHEDULE"
        }
      }
    }
  }

  cluster_addons = {
    coredns    = {}
    kube-proxy = {}
    vpc-cni    = {}
    aws-ebs-csi-driver = {}
  }

  tags = {
    Environment = var.environment
  }
}

module "rds" {
  source  = "terraform-aws-modules/rds/aws"
  version = "6.4.0"

  identifier = "${local.name_prefix}-postgres"

  engine               = "postgres"
  engine_version       = "16.1"
  family               = "postgres16"
  major_engine_version = "16"
  instance_class       = "db.r6g.large"

  allocated_storage     = 200
  max_allocated_storage = 500
  storage_encrypted     = true
  storage_type          = "gp3"

  db_name  = "platform_db"
  username = "platform"
  password = random_password.db_password.result
  port     = 5432

  multi_az               = true
  db_subnet_group_name   = module.vpc.database_subnet_group
  vpc_security_group_ids = [aws_security_group.rds.id]

  backup_window      = "03:00-04:00"
  backup_retention_period = 30
  maintenance_window = "sun:04:00-sun:05:00"

  deletion_protection = true
  skip_final_snapshot = false
  copy_tags_to_snapshot = true

  enabled_cloudwatch_logs_exports = ["postgresql", "upgrade"]
}

module "elasticache" {
  source  = "terraform-aws-modules/elasticache/aws"
  version = "1.2.0"

  cluster_id = "${local.name_prefix}-redis"

  engine               = "redis"
  engine_version       = "7.1"
  node_type            = "cache.r6g.large"
  num_cache_nodes      = 2
  parameter_group_family = "redis7"

  subnet_group_name = module.vpc.elasticache_subnet_group
  security_group_ids = [aws_security_group.redis.id]

  maintenance_window = "sun:05:00-sun:06:00"
  snapshot_window    = "04:00-05:00"
  snapshot_retention_limit = 7
}

module "msk" {
  source  = "terraform-aws-modules/msk-kafka-cluster/aws"
  version = "2.3.0"

  cluster_name = "${local.name_prefix}-kafka"
  kafka_version = "3.5.1"
  number_of_broker_nodes = 3

  broker_node_client_subnets  = module.vpc.private_subnets
  broker_node_security_groups = [aws_security_group.kafka.id]

  broker_node_instance_type = "kafka.m5.large"

  broker_ebs_volume_size = 100

  encryption_in_transit_client_broker = "TLS"
  encryption_in_transit_in_cluster    = true

  configuration_name        = "${local.name_prefix}-kafka-config"
  configuration_server_properties = {
    "auto.create.topics.enable" = "false"
    "default.replication.factor" = "3"
    "min.insync.replicas"       = "2"
    "num.partitions"            = "3"
    "log.retention.hours"       = "168"
    "log.segment.bytes"         = "1073741824"
  }
}

module "s3_documents" {
  source  = "terraform-aws-modules/s3-bucket/aws"
  version = "3.15.1"

  bucket = "${local.name_prefix}-documents"

  control_object_ownership = true
  object_ownership         = "BucketOwnerEnforced"

  versioning = {
    enabled = true
  }

  server_side_encryption_configuration = {
    rule = {
      apply_server_side_encryption_by_default = {
        sse_algorithm = "AES256"
      }
    }
  }

  lifecycle_rule = [
    {
      id      = "expire-old-versions"
      enabled = true
      noncurrent_version_expiration = {
        days = 90
      }
    }
  ]

  # Block public access
  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true
}

resource "random_password" "db_password" {
  length  = 24
  special = false
}

resource "aws_security_group" "rds" {
  name_prefix = "${local.name_prefix}-rds-sg"
  vpc_id      = module.vpc.vpc_id

  ingress {
    from_port = 5432
    to_port   = 5432
    protocol  = "tcp"
    security_groups = [module.eks.cluster_security_group_id]
  }
}

resource "aws_security_group" "redis" {
  name_prefix = "${local.name_prefix}-redis-sg"
  vpc_id      = module.vpc.vpc_id

  ingress {
    from_port = 6379
    to_port   = 6379
    protocol  = "tcp"
    security_groups = [module.eks.cluster_security_group_id]
  }
}

resource "aws_security_group" "kafka" {
  name_prefix = "${local.name_prefix}-kafka-sg"
  vpc_id      = module.vpc.vpc_id

  ingress {
    from_port = 9092
    to_port   = 9094
    protocol  = "tcp"
    security_groups = [module.eks.cluster_security_group_id]
  }
}
