resource "aws_secretsmanager_secret" "db_password" {
  name = "platform-${var.environment}-db-password"
}

resource "aws_secretsmanager_secret_version" "db_password" {
  secret_id     = aws_secretsmanager_secret.db_password.id
  secret_string = random_password.db_password.result
}

resource "aws_secretsmanager_secret" "jwt_secret" {
  name = "platform-${var.environment}-jwt-secret"
}

resource "aws_secretsmanager_secret" "openai_api_key" {
  name = "platform-${var.environment}-openai-api-key"
}

resource "aws_secretsmanager_secret" "smtp_credentials" {
  name = "platform-${var.environment}-smtp-credentials"
}
