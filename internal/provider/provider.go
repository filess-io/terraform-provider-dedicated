package provider

import (
	"fmt"

	"github.com/filess/terraform-provider-dedicated/internal/client"
	"github.com/filess/terraform-provider-dedicated/internal/datasources"
	"github.com/filess/terraform-provider-dedicated/internal/resources"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"api_token": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("FILESS_API_TOKEN", nil),
				Description: "API token for filess.io authentication",
			},
			"api_url": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("FILESS_API_URL", "https://backend.filess.io"),
				Description: "Base URL for filess.io API",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"filess_database": resources.ResourceDatabase(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"filess_engines": datasources.DataSourceEngines(),
			"filess_regions": datasources.DataSourceRegions(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	apiToken := d.Get("api_token").(string)
	apiURL := d.Get("api_url").(string)

	// Validar que el token no esté vacío
	if apiToken == "" {
		return nil, fmt.Errorf("api_token cannot be empty")
	}

	// Validar que la URL no esté vacía
	if apiURL == "" {
		return nil, fmt.Errorf("api_url cannot be empty")
	}

	return client.NewClient(apiURL, apiToken), nil
}
