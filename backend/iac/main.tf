provider "aws" {
    region = "us-east-1"
}

provider "github" {
	token = file("credentials")
	//owner = "cmelgreen"
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

resource "aws_iam_role" "codebuild_iam_role" {
  name = "example"

  assume_role_policy = <<POLICY
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Service": "codebuild.amazonaws.com"
      },
      "Action": "sts:AssumeRole"
    }
  ]
}
POLICY
}

resource "aws_iam_role_policy" "codebuild_role_policy" {
  role = aws_iam_role.codebuild_iam_role.name

  policy = <<POLICY
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Resource": [
        "*"
      ],
      "Action": [
        "logs:CreateLogGroup",
        "logs:CreateLogStream",
        "logs:PutLogEvents"
      ]
    },
    {
      "Effect": "Allow",
      "Action": [
        "s3:*"
      ],
      "Resource": "*"
    }
  ]
}
POLICY
}

locals {
  github_repo = "https://github.com/cmelgreen/personal-site-v2"
}

# resource "aws_codebuild_project" "site_codebuild" {
#   name          = "site-codebuild"
#   description   = "test_site_codebuild_project"
#   build_timeout = "5"
#   service_role  = aws_iam_role.codebuild_iam_role.arn

#   artifacts {
#     type = "S3"
#     name = "."
#     location = aws_s3_bucket.site_bucket.bucket
#     namespace_type = "NONE"
#     packaging = "NONE"
#     encryption_disabled = true
#   }

#   environment {
#     compute_type                = "BUILD_GENERAL1_LARGE"
#     image                       = "aws/codebuild/standard:1.0"
#     type                        = "LINUX_CONTAINER"
#     image_pull_credentials_type = "CODEBUILD"
#   }

#   logs_config {
#     cloudwatch_logs {
#       group_name  = "log-group"
#       stream_name = "log-stream"
#     }
#   }

#   source {
#     type            = "GITHUB"
#     location        = local.github_repo
#     git_clone_depth = 1
#     buildspec = "frontend/buildspec.yml"

#     auth {
#       type = "OAUTH"
#     }
#   }
# }

# resource "aws_codebuild_webhook" "webhook" {
#   project_name = aws_codebuild_project.site_codebuild.name
#   branch_filter = "master"
# }

resource "aws_codebuild_project" "backend_codebuild" {
  name          = "site-backend-codebuild"
  description   = "backend_codebuild_project"
  build_timeout = "5"
  service_role  = aws_iam_role.codebuild_iam_role.arn

  artifacts {
    type = "CODEPIPELINE"
  }

  environment {
    compute_type                = "BUILD_GENERAL1_LARGE"
    image                       = "aws/codebuild/standard:1.0"
    type                        = "LINUX_CONTAINER"
    image_pull_credentials_type = "CODEBUILD"
  }

  logs_config {
    cloudwatch_logs {
      group_name  = "log-group"
      stream_name = "log-stream"
    }
  }

  source {
    type            = "CODEPIPELINE"
    buildspec       = "backend/buildspec.yml"
  }
}

# resource "aws_codebuild_webhook" "backend_webhook" {
#   project_name = aws_codebuild_project.backend_codebuild.name
#   branch_filter = "master"
# }


resource "aws_codedeploy_app" "backend_app" {
    name                    = "backend-app"
}

resource "aws_s3_bucket" "codepipeline_bucket" {
    bucket                  = "cmelgreen-personal-site-pipeline-bucket"
}

resource "aws_codedeploy_deployment_group" "backend_deployment_group" {
    deployment_group_name   = "${aws_codedeploy_app.backend_app.name}-group"
    app_name                = aws_codedeploy_app.backend_app.name

    service_role_arn        = aws_iam_role.codedeploy_iam_role.arn
    autoscaling_groups      = [aws_autoscaling_group.backend_asg.name]

    auto_rollback_configuration {
        enabled             = true
        events              = ["DEPLOYMENT_FAILURE"]
    }
}

resource "aws_autoscaling_group" "backend_asg" {
    name                        = "backend-asg"

    min_size                    = 1
    max_size                    = 1
    desired_capacity            = 1

    health_check_grace_period   = 30
    health_check_type           = "EC2"
    force_delete                = true

    launch_configuration        = aws_launch_configuration.backend_lc.name
    vpc_zone_identifier         = [aws_subnet.public_subnet.id]

    load_balancers              = [aws_elb.backend_elb.name]
}

resource "aws_launch_configuration" "backend_lc" {
    name                        = "backend-lc-${formatdate("YY-MM-DD-HH-mm", timestamp())}"

    image_id                    = "ami-0885b1f6bd170450c"
    instance_type               = "t2.nano"
    user_data                   = "docker run -p 80:80 nginx"

    security_groups             = [aws_security_group.public_http_sg.id]
    iam_instance_profile        = aws_iam_instance_profile.backend_iam_profile.name
    key_name                    = "zoff3"

    associate_public_ip_address = true

    root_block_device {
        volume_type             = "gp2"
        volume_size             = 30
    }

    lifecycle {
        // AWS throws an error if false
        create_before_destroy   = true
    }
}

resource "aws_elb" "backend_elb" {
    name                        = "backend-elb"
    security_groups             = [aws_security_group.public_http_sg.id]
    subnets                     = [aws_subnet.public_subnet.id]

    listener {
        lb_port                 = 80
        lb_protocol             = "HTTP"
        instance_port           = 80
        instance_protocol       = "HTTP"
    }

    health_check {
        healthy_threshold       = 2
        unhealthy_threshold     = 2
        timeout                 = 3
        interval                = 30
        target                  = "HTTP:80/"
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

resource "aws_iam_role" "codedeploy_iam_role" {
    name                        = "codedeploy-iam-role"
    force_detach_policies       = true

    assume_role_policy          = <<POLICY
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "",
      "Effect": "Allow",
      "Principal": {
        "Service": "codedeploy.amazonaws.com"
      },
      "Action": "sts:AssumeRole"
    },
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

resource "aws_iam_instance_profile" "codedeploy_iam_profile" {
    name                        = "codedeploy-iam-profile"
    role                        = aws_iam_role.codedeploy_iam_role.name
}

resource "aws_iam_role_policy_attachment" "codedeploy_iam_policy_attachments" {
    for_each  = toset([
      "arn:aws:iam::aws:policy/service-role/AmazonEC2RoleforAWSCodeDeploy",
      "arn:aws:iam::aws:policy/AWSCodeDeployDeployerAccess",
      "arn:aws:iam::aws:policy/service-role/AWSCodeDeployRole"
    ])

    policy_arn                  = each.value
    role                        = aws_iam_role.codedeploy_iam_role.name
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

############################################################

resource "aws_codepipeline" "codepipeline" {
  name     = "cmelgreen-site-pipeline"
  role_arn = aws_iam_role.codepipeline_role.arn

  artifact_store {
    location = aws_s3_bucket.codepipeline_bucket.bucket
    type     = "S3"
  }

  stage {
    name = "Source"

    action {
      name             = "Source"
      category         = "Source"
      owner            = "AWS"
      provider         = "CodeStarSourceConnection"
      version          = "1"
      output_artifacts = ["source_output"]

      configuration = {
        ConnectionArn    = aws_codestarconnections_connection.github_connection.arn
        FullRepositoryId = "cmelgreen/personal-site-v2"
        BranchName       = "master"
      }

      # configuration = {
      #   Owner             = "cmelgreen"
      #   Repo              = "personal-site-v2"
      #   Branch            = "master"
      # }
    }
  }

  stage {
    name = "Build"

    action {
      name             = "Build"
      category         = "Build"
      owner            = "AWS"
      provider         = "CodeBuild"
      input_artifacts  = ["source_output"]
      output_artifacts = ["build_output"]
      version          = "1"

      configuration = {
        ProjectName = aws_codebuild_project.backend_codebuild.name
      }
    }
  }

  stage {
    name = "Deploy"

    action {
      name            = "Deploy"
      category        = "Deploy"
      owner           = "AWS"
      provider        = "CodeDeploy"
      input_artifacts = ["build_output"]
      version         = "1"

      configuration = {
          ApplicationName = aws_codedeploy_app.backend_app.name
          DeploymentGroupName = aws_codedeploy_deployment_group.backend_deployment_group.deployment_group_name
      }
    }
  }
}

resource "aws_iam_role" "codepipeline_role" {
  name = "codepipeline-role"

  assume_role_policy = <<POLICY
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Service": "codepipeline.amazonaws.com"
      },
      "Action": "sts:AssumeRole"
    }
  ]
}
POLICY
}

resource "aws_iam_role_policy" "codepipeline_policy" {
  name = "codepipeline-policy"
  role = aws_iam_role.codepipeline_role.id

  policy = <<POLICY
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect":"Allow",
      "Action": [
        "s3:*"
      ],
      "Resource": [
        "${aws_s3_bucket.codepipeline_bucket.arn}",
        "${aws_s3_bucket.codepipeline_bucket.arn}/*"
      ]
    },
    {
      "Effect": "Allow",
      "Action": [
        "codebuild:*"
      ],
      "Resource": "*"
    },{
      "Effect": "Allow",
      "Action": "codestar-connections:*",
      "Resource": "*"
        }
  ]
}
POLICY
}

# locals {
#   webhook_secret = "super-secret"
# }

# resource "aws_codepipeline_webhook" "codepipeline_webhook" {
#   name            = "test-webhook-github-bar"
#   authentication  = "GITHUB_HMAC"
#   target_action   = "Source"
#   target_pipeline = aws_codepipeline.codepipeline.name

#   authentication_configuration {
#     secret_token = local.webhook_secret
#   }

#   filter {
#     json_path    = "$.ref"
#     match_equals = "refs/heads/{Branch}"
#   }
# }

resource "aws_codestarconnections_connection" "github_connection" {
  name          = "github-connection"
  provider_type = "GitHub"
}

# # Wire the CodePipeline webhook into a GitHub repository.
# resource "github_repository_webhook" "bar" {
#   repository = local.github_repo

#   configuration {
#     url          = aws_codepipeline_webhook.codepipeline_webhook.url
#     content_type = "json"
#     insecure_ssl = true
#     secret       = local.webhook_secret
#   }

#   events = ["push"]
# }