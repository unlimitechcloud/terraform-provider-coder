package instance

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func InternalSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"result": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: OutputSchema(),
			},
		},
		"store": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}