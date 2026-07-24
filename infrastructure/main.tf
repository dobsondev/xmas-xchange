provider "aws" {
  region = "ca-central-1"
}

resource "aws_s3_bucket" "dobsondev-family-xmas-xchange" {
  bucket              = "dobsondev-family-xmas-xchange"
  force_destroy       = false
  object_lock_enabled = false
  tags                = {}
}

resource "aws_s3_bucket_server_side_encryption_configuration" "this" {
  bucket = aws_s3_bucket.dobsondev-family-xmas-xchange.id

  rule {
    bucket_key_enabled       = true
    blocked_encryption_types = ["SSE-C"]

    apply_server_side_encryption_by_default {
      sse_algorithm = "AES256"
    }
  }
}

resource "aws_s3_bucket_versioning" "this" {
  bucket = aws_s3_bucket.dobsondev-family-xmas-xchange.id

  versioning_configuration {
    status = "Disabled"
  }
}

resource "aws_s3_bucket_public_access_block" "this" {
  bucket = aws_s3_bucket.dobsondev-family-xmas-xchange.id

  block_public_acls       = true
  ignore_public_acls      = true
  block_public_policy     = true
  restrict_public_buckets = true
}

resource "aws_iam_role" "github_actions_xmas_xchange" {
  name = "github-actions-xmas-xchange-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = "sts:AssumeRoleWithWebIdentity"
        Principal = {
          Federated = "arn:aws:iam::157437238176:oidc-provider/token.actions.githubusercontent.com"
        }
        Condition = {
          StringEquals = {
            "token.actions.githubusercontent.com:aud" = "sts.amazonaws.com"
          }
          StringLike = {
            "token.actions.githubusercontent.com:sub" = [
              "repo:dobsondev/xmas-xchange:*"
            ]
          }
        }
      }
    ]
  })
}

resource "aws_iam_role_policy" "xmas_exchange_s3_access" {
  name = "XmasExchangeS3Access"
  role = aws_iam_role.github_actions_xmas_xchange.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect   = "Allow"
        Action   = ["s3:PutObject", "s3:GetObject"]
        Resource = "arn:aws:s3:::dobsondev-family-xmas-xchange/*"
      },
      {
        Effect   = "Allow"
        Action   = "s3:ListBucket"
        Resource = "arn:aws:s3:::dobsondev-family-xmas-xchange"
      }
    ]
  })
}

resource "aws_iam_user" "dobsondev_family_xmas_xchange" {
  name = "dobsondev-family-xmas-xchange"
  path = "/"

  tags = {
    Purpose = "Used by the dobsondev-xmas-xchange Python script to upload the results to S3."
  }
}

resource "aws_iam_user_policy_attachment" "s3_full_access" {
  user       = aws_iam_user.dobsondev_family_xmas_xchange.name
  policy_arn = "arn:aws:iam::aws:policy/AmazonS3FullAccess"
}