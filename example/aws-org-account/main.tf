variable "role_arn" {
  type = string
}

variable "name" {
  type = string
}


variable "email" {
  type = string
}

provider "aws" {
  assume_role {
    role_arn = var.role_arn
  }
}

resource "aws_organizations_account" "this" {
  name  = var.name
  email = var.email
}

output "account" {
  value = aws_organizations_account.this
}