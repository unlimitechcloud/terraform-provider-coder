package volume

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Esquema de outputs anidados, compatible con la transformación a listas de un elemento
func OutputSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}