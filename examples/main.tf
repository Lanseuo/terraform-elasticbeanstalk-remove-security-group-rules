terraform {
  required_providers {
    elasticbeanstalk-remove-security-group-rules = {
      version = "~> 0.1.0"
      source  = "hild.dev/edu/elasticbeanstalk-remove-security-group-rules"
    }
  }
}

provider "aws" {
  region = "eu-central-1"
}

provider "elasticbeanstalk-remove-security-group-rules" {
  region = "eu-central-1"
}

resource "aws_elastic_beanstalk_application" "application" {
  name        = "my-application"
  description = "My application"
}

resource "aws_elastic_beanstalk_environment" "environment" {
  name                = "my-environment"
  application         = aws_elastic_beanstalk_application.application.name
  solution_stack_name = "64bit Amazon Linux 2015.03 v2.0.3 running Go 1.4"
}

resource "elasticbeanstalk-remove-security-group-rules_action" "remove" {
  elasticbeanstalk_environment_id = aws_elastic_beanstalk_environment.environment.id
}
