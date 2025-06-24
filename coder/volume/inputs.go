package volume

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func InputSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"subnet_id": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The subnet id where the EC2 instance will run.",
		},
		"name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Device name (e.g., coder).",
		},
		"size": {
			Type:        schema.TypeInt,
			Required:    true,
			Description: "Volume size (GB).",
		},
		"type": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Volume type (e.g., gp3).",
		},
		"coder": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Coder-specific configuration.",
		},
	}
}
