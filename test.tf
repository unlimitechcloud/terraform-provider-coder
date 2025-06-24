terraform {
  required_providers {
    remote = {
      source  = "unlimitechcloud/remote"
      version = "0.0.0"
    }
  }
}

provider "remote" {
  alias  = "ec2"
  lambda = "CoderResourceProxy"
  region = "us-east-1"
}

resource "coder_volume" "dev" {
  provider  = remote.ec2
  subnet_id = "subnet-083224448dace4c25"
  name      = "home"
  type      = "gp3"
  size      = 50
  coder {
    workspace {
      access_port       = 8080
      access_url        = "https://dev.coder.example.com/myws1"
      id                = "myws2"
      is_prebuild       = false
      is_prebuild_claim = false
      name              = "myws2"
      owner {
        email             = "alice@example.com"
        full_name         = "Alice Example"
        groups            = ["developers"]
        id                = "user-1234"
        login_type        = "password"
        name              = "alice"
        oidc_access_token = ""
        rbac_roles {
          name   = "developer"
          org_id = "org-001"
        }
        session_token   = "SESSION_TOKEN_SAMPLE"
        ssh_private_key = "PRIVATE_KEY_SAMPLE"
        ssh_public_key  = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC7..."
      }
      prebuild_count   = 0
      start_count      = 1
      template_id      = "template-5678"
      template_name    = "ubuntu-dev"
      template_version = "1.0.0"
      transition       = "start"
    }
  }
}

resource "coder_instance" "dev" {
  count              = 1
  provider           = remote.ec2
  require_on_demand  = false
  ami_id             = "ami-020cba7c55df1f615"
  instance_type      = "t3.small"
  subnet_id          = "subnet-083224448dace4c25"
  security_group_ids = ["sg-08695b63da890ed3a"]
  key_pair_name      = "coder"
  user_data          = ""
  root_volume_type   = "gp3"
  root_volume_size   = 20
  home_volume_id     = coder_volume.dev.id
  tags = { }
  coder {
    workspace {
      access_port       = 8080
      access_url        = "https://dev.coder.example.com/myws1"
      id                = "myws2"
      is_prebuild       = false
      is_prebuild_claim = false
      name              = "myws2"
      owner {
        email             = "alice@example.com"
        full_name         = "Alice Example"
        groups            = ["developers"]
        id                = "user-1234"
        login_type        = "password"
        name              = "alice"
        oidc_access_token = ""
        rbac_roles {
          name   = "developer"
          org_id = "org-001"
        }
        session_token   = "SESSION_TOKEN_SAMPLE"
        ssh_private_key = "PRIVATE_KEY_SAMPLE"
        ssh_public_key  = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC7..."
      }
      prebuild_count   = 0
      start_count      = 1
      template_id      = "template-5678"
      template_name    = "ubuntu-dev"
      template_version = "1.0.0"
      transition       = "start"
    }
  }
}

output "volume" {
  value = jsonencode(coder_volume.dev.result[0])
}

output "instance" {
  value = jsonencode(coder_instance.dev[0].result[0])
}

