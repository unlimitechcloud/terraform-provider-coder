package remote

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/unlimitechcloud/terraform-provider-remote/remote/base"
	"github.com/unlimitechcloud/terraform-provider-remote/remote/volume"
	"github.com/unlimitechcloud/terraform-provider-remote/remote/instance"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"lambda": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("REMOTE_LAMBDA", nil),
				Description: "Name or ARN of the Lambda function handling lifecycle.",
			},
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("AWS_REGION", nil),
				Description: "AWS region for the Lambda function if using name instead of ARN.",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"coder_volume": volume.Resource(),
			"coder_instance": instance.Resource(),
		},
		ConfigureContextFunc: configureProvider,
	}
}

func configureProvider(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	lambdaName := d.Get("lambda").(string)
	region := d.Get("region").(string)

	var sess *session.Session
	if strings.HasPrefix(lambdaName, "arn:") {
		sess = session.Must(session.NewSession())
	} else {
		if region == "" {
			return nil, diag.Errorf("region is required when lambda is not an ARN")
		}
		awsCfg := aws.NewConfig().WithRegion(region)
		sess = session.Must(session.NewSession(awsCfg))
	}

	return &base.RemoteClient{
		LambdaName: lambdaName,
		Svc:        lambda.New(sess),
	}, nil
}