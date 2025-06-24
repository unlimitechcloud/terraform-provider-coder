package instance

import (
	"github.com/unlimitechcloud/terraform-provider-remote/remote/base"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Resource() *schema.Resource {
	return base.NewRemoteResource(
		"instance",
		InputSchema(),
		OutputSchema(),
		InternalSchema(),
		func(meta interface{}) *base.RemoteClient {
			// Puedes adaptar esto para tu provider real.
			return meta.(*base.RemoteClient)
		},
		nil, // Puedes pasar handlers custom aquí si quieres lógica especial.
	)
}