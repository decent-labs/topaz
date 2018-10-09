provider "aws" {
  access_key = "${ var.aws-access-key }"
  secret_key = "${ var.aws-secret-key }"
  region     = "${ var.aws-region }"
}

# API

resource "aws_instance" "topaz-api" {
  ami = "ami-8c122be9"
  instance_type = "t2.micro"
  key_name = "topaz-api-ssh"

  connection {
    type = "ssh"
    user = "ec2-user"
    private_key = "${ file("~/.ssh/id_rsa") }"
  }

  security_groups = [
    "default"
  ]

  provisioner "remote-exec" {
    inline = [
      "mkdir ~/go",
      "mkdir ~/go/src",
      "mkdir ~/go/src/topaz.io",
      "mkdir ~/go/src/topaz.io/topaz"
    ]
  }

  provisioner "file" {
    source = "../../main.go"
    destination = "~/go/src/topaz.io/topaz/main.go"
  }

  provisioner "remote-exec" {
    script = "./install-api.sh"
  }

  tags {
    Name = "topaz-api"
  }
}

resource "aws_key_pair" "topaz-api-ssh" {
  key_name = "topaz-api-ssh"
  public_key = "${ file("~/.ssh/id_rsa.pub") }"
}
