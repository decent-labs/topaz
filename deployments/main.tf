provider "aws" {
  access_key = "${ var.aws-access-key }"
  secret_key = "${ var.aws-secret-key }"
  region     = "${ var.aws-region }"
}

# batch

resource "aws_instance" "topaz-batch" {
  ami = "ami-8c122be9"
  instance_type = "t2.micro"
  key_name = "topaz-batch-ssh"

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
      "mkdir ~/go/src/topaz.io/batch"
    ]
  }

  provisioner "file" {
    source = "../../main.go"
    destination = "~/go/src/topaz.io/batch/main.go"
  }

  provisioner "remote-exec" {
    script = "./install-batch.sh"
  }

  tags {
    Name = "topaz-batch"
  }
}

resource "aws_key_pair" "topaz-batch-ssh" {
  key_name = "topaz-batch-ssh"
  public_key = "${ file("~/.ssh/id_rsa.pub") }"
}

# collect

resource "aws_instance" "topaz-collect" {
  ami = "ami-8c122be9"
  instance_type = "t2.micro"
  key_name = "topaz-collect-ssh"

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
      "mkdir ~/go/src/topaz.io/collect"
    ]
  }

  provisioner "file" {
    source = "../../collect/main.go"
    destination = "~/go/src/topaz.io/collect/main.go"
  }

  provisioner "remote-exec" {
    script = "./install-collect.sh"
  }

  tags {
    Name = "topaz-collect"
  }
}

resource "aws_key_pair" "topaz-collect-ssh" {
  key_name = "topaz-collect-ssh"
  public_key = "${ file("~/.ssh/id_rsa.pub") }"
}

# ipfs

resource "aws_instance" "topaz-ipfs" {
  ami = "ami-8c122be9"
  instance_type = "t2.micro"
  key_name = "topaz-ipfs-ssh"

  connection {
    type = "ssh"
    user = "ec2-user"
    private_key = "${ file("~/.ssh/id_rsa") }"
  }

  security_groups = [
    "${ aws_security_group.topaz-public.name }"
  ]

  provisioner "remote-exec" {
    script = "./install-ipfs.sh"
  }

  tags {
    Name = "topaz-ipfs"
  }
}

resource "aws_key_pair" "topaz-ipfs-ssh" {
  key_name = "topaz-ipfs-ssh"
  public_key = "${ file("~/.ssh/id_rsa.pub") }"
}

resource "aws_security_group" "topaz-public" {
  name = "topaz-public"
  description = "public for development"

  ingress {
    from_port   = 0
    to_port     = 65535
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags {
    Name = "topaz-public"
  }
}
