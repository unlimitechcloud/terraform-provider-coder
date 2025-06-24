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