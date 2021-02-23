provider "aws" {
    region = "us-east-1"
}

resource "aws_s3_bucket" "site_bucket" {
  bucket = "cm-personal-site-bucket"
  acl    = "public-read"
  policy = <<POLICY
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Sid": "PublicReadGetObject",
            "Effect": "Allow",
            "Principal": "*",
            "Action": "s3:GetObject",
            "Resource": "arn:aws:s3:::cm-personal-site-bucket/*"
        }
    ]
}
POLICY

  website {
    index_document = "index.html"
    error_document = "error.html"
  }

  tags = {
    Name        = "test buck"
  }
}

locals {
  s3_origin_id = "s3-${aws_s3_bucket.site_bucket.bucket}"
}

resource "aws_cloudfront_distribution" "site_distribution" {
  origin {
    domain_name = aws_s3_bucket.site_bucket.bucket_regional_domain_name
    origin_id   = local.s3_origin_id
  }

  enabled             = true
  is_ipv6_enabled     = true
  comment             = "test"
  default_root_object = "index.html"


  default_cache_behavior {
    allowed_methods  = ["GET", "HEAD"]
    cached_methods   = ["GET", "HEAD"]
    target_origin_id = local.s3_origin_id

    forwarded_values {
      query_string = false

      cookies {
        forward = "none"
      }
    }

    viewer_protocol_policy = "allow-all"
    min_ttl                = 0
    default_ttl            = 3600
    max_ttl                = 86400
  }

  # Cache behavior with precedence 0
  ordered_cache_behavior {
    path_pattern     = "/content/immutable/*"
    allowed_methods  = ["GET", "HEAD", "OPTIONS"]
    cached_methods   = ["GET", "HEAD", "OPTIONS"]
    target_origin_id = local.s3_origin_id

    forwarded_values {
      query_string = false
      headers      = ["Origin"]

      cookies {
        forward = "none"
      }
    }

    min_ttl                = 0
    default_ttl            = 86400
    max_ttl                = 31536000
    compress               = true
    viewer_protocol_policy = "redirect-to-https"
  }

  # Cache behavior with precedence 1
  ordered_cache_behavior {
    path_pattern     = "/content/*"
    allowed_methods  = ["GET", "HEAD", "OPTIONS"]
    cached_methods   = ["GET", "HEAD"]
    target_origin_id = local.s3_origin_id

    forwarded_values {
      query_string = false

      cookies {
        forward = "none"
      }
    }

    min_ttl                = 0
    default_ttl            = 3600
    max_ttl                = 86400
    compress               = true
    viewer_protocol_policy = "redirect-to-https"
  }

  price_class = "PriceClass_100"

  restrictions {
    geo_restriction {
      restriction_type = "whitelist"
      locations        = ["US", "CA"]
    }
  }

  viewer_certificate {
    cloudfront_default_certificate = true
  }
}

resource "aws_instance" "backend_ec2" {
    ami                         = "ami-0be2609ba883822ec"
    instance_type               = "t2.nano"
    associate_public_ip_address = true
    key_name                    = "zoff3"

    iam_instance_profile        = aws_iam_instance_profile.backend_iam_profile.name
    
    subnet_id                   = aws_subnet.public_subnet.id
    vpc_security_group_ids      = [aws_security_group.public_http_sg.id]

    user_data                   = file("lc_user_data.sh")

    tags = {
      codedeploy = "personal-site-backend"
    }

}

resource "aws_iam_role" "backend_iam_role" {
    name                        = "backend-iam-role"
    force_detach_policies       = true

    assume_role_policy          = <<POLICY
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "",
      "Effect": "Allow",
      "Principal": {
          "Service": [
          "ec2.amazonaws.com"
          ]
      },
      "Action": "sts:AssumeRole"
    }
  ]
}
POLICY
}

resource "aws_iam_instance_profile" "backend_iam_profile" {
    name                        = "backend-iam-profile"
    role                        = aws_iam_role.backend_iam_role.name
}

resource "aws_iam_role_policy_attachment" "backend_iam_policy_attachments" {
    for_each  = toset([
      "arn:aws:iam::aws:policy/service-role/AmazonEC2RoleforAWSCodeDeploy",
      "arn:aws:iam::aws:policy/AmazonSSMReadOnlyAccess"
    ])

    policy_arn                  = each.value
    role                        = aws_iam_role.backend_iam_role.name
}

resource "aws_vpc" "vpc" {
    cidr_block              = "10.0.0.0/16"
    enable_dns_support      = true
    enable_dns_hostnames    = true

    tags = {
        Name = "ctx-defs-aas-vpc"
    }
}

resource "aws_internet_gateway" "igw" {
    vpc_id          = aws_vpc.vpc.id
}

resource "aws_subnet" "public_subnet" {
    vpc_id                  = aws_vpc.vpc.id
    cidr_block              = "10.0.1.0/24"
    map_public_ip_on_launch = true
    availability_zone       = "us-east-1a"
}

resource "aws_route_table" "public_rtb" {
    vpc_id          = aws_vpc.vpc.id

    route {
        cidr_block  = "0.0.0.0/0"
        gateway_id  = aws_internet_gateway.igw.id
    }
}

resource "aws_route_table_association" "public_route_assosciation" {
    route_table_id  = aws_route_table.public_rtb.id
    subnet_id       = aws_subnet.public_subnet.id
}

resource "aws_security_group" "public_http_sg" {
	name        = "public_http_sg"

    vpc_id = aws_vpc.vpc.id

    ingress {
        from_port   = 22
        to_port     = 22
        protocol    = "tcp"
        cidr_blocks = ["0.0.0.0/0"]
    }

	ingress {
		from_port   = 80
		to_port     = 80
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