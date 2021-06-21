package zoom

import (
	"terraform-provider-zoom/client"
	tkn "terraform-provider-zoom/token"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"zoom_api_secret": &schema.Schema{
				Type:      schema.TypeString,
				Required:  true,
				DefaultFunc: schema.EnvDefaultFunc("ZOOM_API_SECRET", ""),
			},
			"zoom_api_key": &schema.Schema{
				Type:      schema.TypeString,
				Required:  true,
				DefaultFunc: schema.EnvDefaultFunc("ZOOM_API_KEY", ""),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"zoom_user": resourceUser(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"zoom_user": dataSourceUser(),
		},
		ConfigureFunc:  providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	token := tkn.TokenGenerate(d.Get("zoom_api_secret").(string),d.Get("zoom_api_key").(string))
	return client.NewClient(token), nil
}