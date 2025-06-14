resource "aws_ecr_repository" "newsletter_go" {
  name                 = "newsletter-go"
  image_tag_mutability = "MUTABLE"

  image_scanning_configuration {
    scan_on_push = true
  }
}

output "repository_url" {
  value = aws_ecr_repository.newsletter_go.repository_url
}
