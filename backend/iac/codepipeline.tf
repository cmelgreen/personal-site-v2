provider "github" {
	token = file("credentials")
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

resource "aws_codebuild_project" "site_codebuild" {
  name          = "site-codebuild"
  description   = "test_site_codebuild_project"
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
    buildspec       = "frontend/buildspec.yml"
  }

  cache {
    type    = "LOCAL"
    modes   = ["LOCAL_SOURCE_CACHE"]
  }
  
}

resource "aws_codebuild_project" "backend_codebuild" {
  name          = "site-backend-codebuild"
  description   = "backend_codebuild_project"
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

  cache {
    type    = "LOCAL"
    modes   = ["LOCAL_SOURCE_CACHE"]
  }
}

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
    ec2_tag_filter {
      type = "VALUE_ONLY"
      value = "personal-site-backend"
    }

    auto_rollback_configuration {
        enabled             = true
        events              = ["DEPLOYMENT_FAILURE"]
    }
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
      name             = "Build-Backend"
      category         = "Build"
      owner            = "AWS"
      provider         = "CodeBuild"
      input_artifacts  = ["source_output"]
      output_artifacts = ["backend_build_output"]
      version          = "1"

      configuration = {
        ProjectName = aws_codebuild_project.backend_codebuild.name
      }
    }

    action {
      name             = "Build-Frontend"
      category         = "Build"
      owner            = "AWS"
      provider         = "CodeBuild"
      input_artifacts  = ["source_output"]
      output_artifacts = ["frontend_build_output"]
      version          = "1"

      configuration = {
        ProjectName = aws_codebuild_project.site_codebuild.name
      }
    }
  }

  stage {
    name = "Deploy"

    action {
      name            = "Deploy-Backend"
      category        = "Deploy"
      owner           = "AWS"
      provider        = "CodeDeploy"
      input_artifacts = ["backend_build_output"]
      version         = "1"

      configuration = {
          ApplicationName = aws_codedeploy_app.backend_app.name
          DeploymentGroupName = aws_codedeploy_deployment_group.backend_deployment_group.deployment_group_name
      }
    }

    action {
      name            = "Deploy-Frontend"
      category        = "Deploy"
      owner           = "AWS"
      provider        = "S3"
      input_artifacts = ["frontend_build_output"]
      version         = "1"

      configuration = {
        BucketName = aws_s3_bucket.site_bucket.bucket
        Extract = "true"
      }
    }

      # artifacts {
  #   type = "S3"
  #   name = "."
  #   location = aws_s3_bucket.site_bucket.bucket
  #   namespace_type = "NONE"
  #   packaging = "NONE"
  #   encryption_disabled = true
  # }


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
      "Resource": "*"
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
    },{
      "Effect": "Allow",
      "Action": [
        "codedeploy:*"
      ],
      "Resource": "*"
    }
  ]
}
POLICY
}

resource "aws_codestarconnections_connection" "github_connection" {
  name          = "github-connection"
  provider_type = "GitHub"
}