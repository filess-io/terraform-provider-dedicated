package datasources

import (
	"context"
	"fmt"

	"github.com/filess/terraform-provider-dedicated/internal/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceRegions() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRegionsRead,
		Schema: map[string]*schema.Schema{
			"regions": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of available regions",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Region ID",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Region name",
						},
						"region_code": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Region code",
						},
						"availability_domain": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Availability domain",
						},
					},
				},
			},
		},
	}
}

func dataSourceRegionsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)

	resp, err := c.Get("/api/v1/regions")
	if err != nil {
		return diag.FromErr(err)
	}

	regions := resp.Data.([]interface{})

	regionList := make([]map[string]interface{}, len(regions))
	for i, region := range regions {
		r := region.(map[string]interface{})

		// Manejar id que puede ser string o float64
		var idStr string
		switch v := r["id"].(type) {
		case float64:
			idStr = fmt.Sprintf("%.0f", v)
		case string:
			idStr = v
		default:
			idStr = fmt.Sprintf("%v", v)
		}

		regionList[i] = map[string]interface{}{
			"id":                  idStr,
			"name":                r["name"],
			"region_code":         r["regionCode"],
			"availability_domain": r["availabilityDomain"],
		}
	}

	d.SetId("regions")
	d.Set("regions", regionList)

	return nil
}
