variable "region" {
  type = string
  default = "us-west-1"
}

variable "instance_type" {
  type = string
  default = "t3a.small"
}

variable "key_name" {
  type = string
  description = "The ssh key pair to connect to the ec2 instances"
}

variable "security_group_id" {
  type = string
  description = "AWS security group for the ec2 instances"
}

variable "iam_instance_profile" {
  type = string
  description = "AWS IAM instance profile for the ec2 instances"
}

variable "eip_red" {
  type = string
  description = "Elastic ip for the tf2 server 'red'"
}


