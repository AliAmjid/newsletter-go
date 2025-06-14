resource "aws_security_group" "newsletter_ec2_sg" {
  name        = "newsletter-ec2-sg"
  description = "Allow SSH and HTTP access"

  ingress {
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    from_port   = 80
    to_port     = 3000
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

resource "aws_instance" "newsletter_ec2" {
  ami                         = var.ec2_ami_id
  instance_type               = var.ec2_instance_type
  vpc_security_group_ids      = [aws_security_group.newsletter_ec2_sg.id]
  associate_public_ip_address = true
  key_name                    = var.ec2_key_name

  user_data = <<-EOF
    #!/bin/bash
    set -e
    # Install Docker
    if ! command -v docker &> /dev/null; then
      apt-get update
      apt-get install -y docker.io
      usermod -aG docker ubuntu
      systemctl enable docker
      systemctl start docker
    fi
    # Install AWS CLI v2
    if ! command -v aws &> /dev/null; then
      apt-get install -y unzip curl
      curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "awscliv2.zip"
      unzip awscliv2.zip
      ./aws/install
    fi
  EOF

  tags = {
    Name = "newsletter-ec2"
  }
}

output "ec2_public_ip" {
  value = aws_instance.newsletter_ec2.public_ip
}
