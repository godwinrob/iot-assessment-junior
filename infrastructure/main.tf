terraform {
  required_version = ">= 0.12.24"
}

provider aws {
  version = ">= 2.57.0"
  region = "us-east-1"
}

# Your HCL goes below! You got this!

resource "aws_dynamodb_table" "ddbtable" {
  name             = "Users"
  hash_key         = "email"
  billing_mode   = "PROVISIONED"
  read_capacity  = 5
  write_capacity = 5
  attribute {
    name = "email"
    type = "S"
  }
}

resource "aws_iam_role_policy" "lambda_policy" {
  name = "lambda_policy"
  role = aws_iam_role.role_for_LDC.id

  policy = file("policy.json")
}