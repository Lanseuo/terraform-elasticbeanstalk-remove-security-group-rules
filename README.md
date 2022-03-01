# Terraform Elasticbeanstalk Remove Security Group Rules

**Terraform Elasticbeanstalk Remove Security Group Rules** is a Terraform provider that deletes all ingress rules from the default security group of an Elastic Beanstalk environment.

By default, Elasticbeanstalk creates a security group for every environment that allows ingress traffic on port 80 and/or 443. For security purposes, it makes sense in certain situation to remove these security group rules and attach custom security groups. However, Elastic Beanstalk doesn't allow overwriting this behavior ([described in several issues](https://github.com/hashicorp/terraform-provider-aws/issues/2002)). This Terraform provider provides a workaround for this issue by deleting all security group rules on the default security group that is attached to all EC2 instances of the Elastic Beanstalk environment. This does not impact the security group which is attached to the load balancer of the environment.

> This resource does not represent a real-world entity in AWS, therefore changing or deleting this resource on its own has no immediate effect.

## Quickstart

```terraform
terraform {
  required_providers {
    elasticbeanstalk-remove-security-group-rules = {
      version = "~> 0.1.0"
      source  = "hild.dev/edu/elasticbeanstalk-remove-security-group-rules"
    }
  }
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
```
