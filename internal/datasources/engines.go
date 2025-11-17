package datasources

import (
	"context"
	"fmt"

	"github.com/filess/terraform-provider-dedicated/internal/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceEngines() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceEnginesRead,
		Schema: map[string]*schema.Schema{
			"engines": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of available database engines",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Engine ID",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Engine name",
						},
						"version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Engine version",
						},
						"slug": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Engine slug",
						},
						"active": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether the engine is active",
						},
					},
				},
			},
		},
	}
}

func dataSourceEnginesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)

	resp, err := c.Get("/api/v1/engines")
	if err != nil {
		return diag.FromErr(err)
	}

	engines := resp.Data.([]interface{})

	engineList := make([]map[string]interface{}, len(engines))
	for i, engine := range engines {
		e := engine.(map[string]interface{})

		// Manejar id que puede ser string o float64
		var idStr string
		switch v := e["id"].(type) {
		case float64:
			idStr = fmt.Sprintf("%.0f", v)
		case string:
			idStr = v
		default:
			idStr = fmt.Sprintf("%v", v)
		}

		engineList[i] = map[string]interface{}{
			"id":      idStr,
			"name":    e["name"],
			"version": e["version"],
			"slug":    e["slug"],
			"active":  e["active"],
		}
	}

	d.SetId("engines")
	d.Set("engines", engineList)

	return nil
}
