resource "aws_db_instance" "rds" {
      identifier                = var.IDENTIFIER
      username                  = var.USERNAME
      password                  = var.PASSWORD
      final_snapshot_identifier = "${var.IDENTIFIER}-${var.SNAPSHOT_NAME}"
      skip_final_snapshot       = var.SNAPSHOT_SKIP
      allocated_storage         = 5
      storage_type              = "gp2"
      instance_class            = "db.t2.micro"
      engine                    = "postgres"
      engine_version            = "12.5"
      publicly_accessible       = true

    vpc_security_group_ids    = [aws_security_group.sg.id]
    db_subnet_group_name      = aws_db_subnet_group.subnet_group.name
}

resource "aws_db_subnet_group" "subnet_group" {
    subnet_ids    = [aws_subnet.public_subnet.id]
}

resource "aws_security_group" "rds_sg" {
    vpc_id              = aws_vpc.vpc.id

    ingress {
        from_port       = 5432
        to_port         = 5432
        protocol        = "tcp"
        security_groups = var.INGRESS_SGS
    }

    egress {
        from_port   = 0
        to_port     = 0
        protocol    = "-1"
        cidr_blocks = ["0.0.0.0/0"]
    }
}

resource "aws_ssm_parameter" "user" {
    name        = "/${var.IDENTIFIER}/user"
    type        = "SecureString"
    value       = var.USERNAME
}

resource "aws_ssm_parameter" "password" {
    name        = "/${var.IDENTIFIER}/password"
    type        = "SecureString"
    value       = var.PASSWORD
}

resource "aws_ssm_parameter" "port" {
    name        = "/${var.IDENTIFIER}/port"
    type        = "SecureString"
    value       = var.PORT
}

resource "aws_ssm_parameter" "host" {
    name        = "/${var.IDENTIFIER}/host"
    type        = "SecureString"
    value       = aws_db_instance.rds.address
}

resource "aws_ssm_parameter" "database" {
    name        = "/${var.IDENTIFIER}/database"
    type        = "SecureString"
    value       = var.ENGINE
}