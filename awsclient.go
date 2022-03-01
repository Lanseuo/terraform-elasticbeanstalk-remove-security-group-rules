package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func findElasticbeanstalkSecurityGroups(svc *ec2.EC2, elasticbeanstalkEnvironmentID string) ([]*ec2.SecurityGroup, error) {
	result, err := svc.DescribeSecurityGroups(&ec2.DescribeSecurityGroupsInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("tag:elasticbeanstalk:environment-id"),
				Values: []*string{aws.String(elasticbeanstalkEnvironmentID)},
			},
			{
				Name:   aws.String("tag:aws:cloudformation:logical-id"),
				Values: []*string{aws.String("AWSEBSecurityGroup")},
			},
		},
	})
	if err != nil {
		return []*ec2.SecurityGroup{}, err
	}

	return result.SecurityGroups, nil
}

func removeSecurityGroupIngress(svc *ec2.EC2, securityGroup *ec2.SecurityGroup) error {
	hasIngressRules := len(securityGroup.IpPermissions) > 0
	if !hasIngressRules {
		return nil
	}

	_, err := svc.RevokeSecurityGroupIngress(&ec2.RevokeSecurityGroupIngressInput{
		GroupId:       securityGroup.GroupId,
		IpPermissions: securityGroup.IpPermissions,
	})
	return err
}
