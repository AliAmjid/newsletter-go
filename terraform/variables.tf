variable "aws_region" {
  description = "AWS region to deploy resources"
  type        = string
  default     = "us-central-1"
}

variable "aws_account_id" {
  description = "AWS Account ID"
  type        = string
}

variable "aws_profile" {
  description = "AWS CLI profile to use (from your ~/.aws/credentials)"
  default     = "development"
}
