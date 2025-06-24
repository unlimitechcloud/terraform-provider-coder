package instance

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func InputSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"require_on_demand": {
			Type:        schema.TypeBool,
			Required:    true,
			Description: "Whether the instance should be provisioned on demand.",
		},
		"ami_id": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The ID of the AMI to use.",
		},
		"instance_type": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The EC2 instance type.",
			
		},
		"subnet_id": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The subnet ID where the instance will be launched.",
		},
		"security_group_ids": {
			Type:        schema.TypeList,
			Required:    true,
			Description: "List of security group IDs.",
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"key_pair_name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The name of the SSH key pair.",
		},
		"user_data": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "User data for the instance.",
			Default:     "",
		},
		"root_volume_type": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Root volume type.",
		},
		"root_volume_size": {
			Type:        schema.TypeInt,
			Required:    true,
			Description: "Root volume size (GB).",
		},
		"home_volume_id": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Home volume id.",
		},
		"tags": {
			Type:        schema.TypeMap,
			Optional:    true,
			Description: "Instance tags.",
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"coder": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Coder-specific configuration.",
		},
	}
}