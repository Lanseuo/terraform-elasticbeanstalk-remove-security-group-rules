package main

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Required: true,
			},
			"profile": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"elasticbeanstalk-remove-security-group-rules_action": resource(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	var diags diag.Diagnostics

	region := d.Get("region").(string)
	profile := d.Get("profile").(string)

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region)},
	)
	if err != nil {
		return nil, diag.FromErr(err)
	}

	if profile != "" {
		sess, err = session.NewSessionWithOptions(session.Options{
			Config: aws.Config{
				Region: aws.String(region),
			},
			SharedConfigState: session.SharedConfigEnable,
			Profile:           profile,
		})
		if err != nil {
			return nil, diag.FromErr(err)
		}
	}

	_, err = sess.Config.Credentials.Get()
	if err != nil {
		return nil, diag.FromErr(err)
	}

	svc := ec2.New(sess)

	return svc, diags
}

func resource() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCreate,
		ReadContext:   resourceRead,
		UpdateContext: resourceUpdate,
		DeleteContext: resourceDelete,
		Schema: map[string]*schema.Schema{
			"elasticbeanstalk_environment_id": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	elasticbeanstalkEnvironmentID := d.Get("elasticbeanstalk_environment_id").(string)

	var diags diag.Diagnostics

	svc := m.(*ec2.EC2)

	securityGroups, err := findElasticbeanstalkSecurityGroups(svc, elasticbeanstalkEnvironmentID)
	if err != nil {
		return diag.FromErr(err)
	}

	if len(securityGroups) == 0 {
		return diag.Errorf("Unable to find security group that is attached to Elasticbeanstalk environment %s", elasticbeanstalkEnvironmentID)
	}
	if len(securityGroups) > 1 {
		return diag.Errorf("Found several security groups that are attached to Elasticbeanstalk environment %s", elasticbeanstalkEnvironmentID)
	}

	err = removeSecurityGroupIngress(svc, securityGroups[0])
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(elasticbeanstalkEnvironmentID)

	return diags
}

func resourceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	elasticbeanstalkEnvironmentID := d.Get("elasticbeanstalk_environment_id").(string)

	var diags diag.Diagnostics

	svc := m.(*ec2.EC2)

	securityGroups, err := findElasticbeanstalkSecurityGroups(svc, elasticbeanstalkEnvironmentID)
	if err != nil {
		return diag.FromErr(err)
	}

	if len(securityGroups) == 0 {
		return diag.Errorf("Unable to find security group that is attached to Elasticbeanstalk environment %s", elasticbeanstalkEnvironmentID)
	}
	if len(securityGroups) > 1 {
		return diag.Errorf("Found several security groups that are attached to Elasticbeanstalk environment %s", elasticbeanstalkEnvironmentID)
	}

	securityGroup := securityGroups[0]
	hasIngressRules := len(securityGroup.IpPermissions) > 0
	if !hasIngressRules {
		d.SetId(elasticbeanstalkEnvironmentID)
	} else {
		d.SetId("")
	}

	return diags
}

func resourceUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return resourceRead(ctx, d, m)
}

func resourceDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	diags = append(diags, diag.Diagnostic{
		Severity: diag.Warning,
		Summary:  "Deleting this resource on its own has no immediate effect.",
		Detail:   "This resource does not represent a real-world entity in AWS, therefore changing or deleting this resource on its own has no immediate effect.",
	})

	d.SetId("")

	return diags
}
