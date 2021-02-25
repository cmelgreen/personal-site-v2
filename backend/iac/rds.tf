resource "aws_db_instance" "rds" {
      identifier                = "personal-site-db"
      username                  = "postgres"
      password                  = "postgres-personal-site"
      final_snapshot_identifier = "personal-site-db-snahpshot"
      skip_final_snapshot       = true
      allocated_storage         = 5
      storage_type              = "gp2"
      instance_class            = "db.t2.micro"
      engine                    = "postgres"
      engine_version            = "12.5"
      publicly_accessible       = true

    vpc_security_group_ids    = [aws_security_group.rds_sg.id]
    db_subnet_group_name      = aws_db_subnet_group.subnet_group.name
}

resource "aws_db_subnet_group" "subnet_group" {
    subnet_ids    = [
        aws_subnet.public_subnet.id,
        aws_subnet.backup_subnet.id
    ]
}

resource "aws_security_group" "rds_sg" {
    vpc_id              = aws_vpc.vpc.id

    ingress {
        from_port       = 5432
        to_port         = 5432
        protocol        = "tcp"
        security_groups = [aws_security_group.public_http_sg.id]
        
        cidr_blocks = [file("local_ip_credentials")]
    }

    egress {
        from_port   = 0
        to_port     = 0
        protocol    = "-1"
        cidr_blocks = ["0.0.0.0/0"]
    }
}

resource "aws_ssm_parameter" "user" {
    name        = "/personal-site-db/user"
    type        = "SecureString"
    value       = "postgres"
}

resource "aws_ssm_parameter" "password" {
    name        = "/personal-site-db/password"
    type        = "SecureString"
    value       = "postgres"
}

resource "aws_ssm_parameter" "port" {
    name        = "/personal-site-db/port"
    type        = "SecureString"
    value       = 5432
}

resource "aws_ssm_parameter" "host" {
    name        = "/personal-site-db/host"
    type        = "SecureString"
    value       = aws_db_instance.rds.address
}

resource "aws_ssm_parameter" "database" {
    name        = "/personal-site-db/database"
    type        = "SecureString"
    value       = "postgres"
}

resource "aws_subnet" "backup_subnet" {
    vpc_id                  = aws_vpc.vpc.id
    cidr_block              = "10.0.2.0/24"
    availability_zone       = "us-east-1b"
}

resource "aws_route_table" "backup_rtb" {
    vpc_id          = aws_vpc.vpc.id

    route {
        cidr_block  = "0.0.0.0/0"
        gateway_id  = aws_internet_gateway.igw.id
    }
}

resource "aws_route_table_association" "backup_subnet" {
    subnet_id      = aws_subnet.backup_subnet.id
    route_table_id = aws_route_table.backup_rtb.id
}