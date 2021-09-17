package zoom

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"terraform-provider-zoom/client"
	"terraform-provider-zoom/token"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"zoom_api_secret": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("ZOOM_API_SECRET", ""),
			},
			"zoom_api_key": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("ZOOM_API_KEY", ""),
			},
			"zoom_timeout_minutes": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ZOOM_TIMEOUT_MINUTES", 2),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"zoom_user": resourceUser(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"zoom_user": dataSourceUser(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	var api_secret string
	if v, ok := d.GetOk("zoom_api_secret"); ok {
		api_secret = v.(string)
	}
	var api_key string
	if v, ok := d.GetOk("zoom_api_key"); ok {
		api_key = v.(string)
	}
	var timeout_minutes int
	if v, ok := d.GetOk("zoom_timeout_minutes"); ok {
		timeout_minutes = v.(int)
	}
	apiToken := token.GenerateToken(api_secret, api_key)
	return client.NewClient(apiToken, timeout_minutes), nil
}
