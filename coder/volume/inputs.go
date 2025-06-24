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
			Type:        schema.TypeList,
			Required:    true,
			Description: "Coder-specific configuration.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"workspace": {
						Type:        schema.TypeList,
						Required:    true,
						Description: "Coder workspace configuration.",
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"access_port": {
									Type:        schema.TypeInt,
									Required:    true,
									Description: "Workspace access port.",
								},
								"access_url": {
									Type:        schema.TypeString,
									Required:    true,
									Description: "Workspace access URL.",
								},
								"id": {
									Type:        schema.TypeString,
									Required:    true,
									Description: "Workspace ID.",
								},
								"is_prebuild": {
									Type:        schema.TypeBool,
									Required:    true,
									Description: "Is workspace prebuilt.",
								},
								"is_prebuild_claim": {
									Type:        schema.TypeBool,
									Required:    true,
									Description: "Is workspace claimed as prebuild.",
								},
								"name": {
									Type:        schema.TypeString,
									Required:    true,
									Description: "Workspace name.",
								},
								"owner": {
									Type:        schema.TypeList,
									Required:    true,
									Description: "Workspace owner information.",
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"email": {
												Type:        schema.TypeString,
												Required:    true,
												Description: "Owner email.",
											},
											"full_name": {
												Type:        schema.TypeString,
												Required:    true,
												Description: "Owner full name.",
											},
											"groups": {
												Type:        schema.TypeList,
												Required:    true,
												Description: "Owner groups.",
												Elem:        &schema.Schema{Type: schema.TypeString},
											},
											"id": {
												Type:        schema.TypeString,
												Required:    true,
												Description: "Owner ID.",
											},
											"login_type": {
												Type:        schema.TypeString,
												Required:    true,
												Description: "Owner login type.",
											},
											"name": {
												Type:        schema.TypeString,
												Required:    true,
												Description: "Owner name.",
											},
											"oidc_access_token": {
												Type:        schema.TypeString,
												Optional:    true,
												Default:     "",
												Description: "Owner OIDC access token.",
											},
											"rbac_roles": {
												Type:        schema.TypeList,
												Required:    true,
												Description: "RBAC roles for the owner.",
												Elem: &schema.Resource{
													Schema: map[string]*schema.Schema{
														"name": {
															Type:        schema.TypeString,
															Required:    true,
															Description: "Role name.",
														},
														"org_id": {
															Type:        schema.TypeString,
															Required:    true,
															Description: "Organization ID.",
														},
													},
												},
											},
											"session_token": {
												Type:        schema.TypeString,
												Required:    true,
												Description: "Owner session token.",
											},
											"ssh_private_key": {
												Type:        schema.TypeString,
												Required:    true,
												Description: "Owner SSH private key.",
											},
											"ssh_public_key": {
												Type:        schema.TypeString,
												Required:    true,
												Description: "Owner SSH public key.",
											},
										},
									},
								},
								"prebuild_count": {
									Type:        schema.TypeInt,
									Required:    true,
									Description: "Workspace prebuild count.",
								},
								"start_count": {
									Type:        schema.TypeInt,
									Required:    true,
									Description: "Workspace start count.",
								},
								"template_id": {
									Type:        schema.TypeString,
									Required:    true,
									Description: "Workspace template ID.",
								},
								"template_name": {
									Type:        schema.TypeString,
									Required:    true,
									Description: "Workspace template name.",
								},
								"template_version": {
									Type:        schema.TypeString,
									Required:    true,
									Description: "Workspace template version.",
								},
								"transition": {
									Type:        schema.TypeString,
									Required:    true,
									Description: "Workspace transition action.",
								},
							},
						},
					},
				},
			},
		},
	}
}
