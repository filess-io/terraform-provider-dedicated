package resources

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/filess/terraform-provider-filess/internal/client"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceDatabase() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDatabaseCreate,
		ReadContext:   resourceDatabaseRead,
		UpdateContext: resourceDatabaseUpdate,
		DeleteContext: resourceDatabaseDelete,
		Schema: map[string]*schema.Schema{
			"organization_slug": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Organization slug",
			},
			"namespace_slug": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Namespace slug",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Database name",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Database description",
			},
			"engine_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Database engine ID",
			},
			"region_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Region ID",
			},
			"database_plan": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Description: "Database plan configuration",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"billable_items": {
							Type:        schema.TypeSet,
							Required:    true,
							Description: "Set of billable items",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"billable_item_id": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Billable item ID",
									},
									"quantity": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Quantity of the billable item",
									},
								},
							},
						},
					},
				},
			},
			"ip_whitelist_ids": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List of IP whitelist IDs",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"ssh_key_ids": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List of SSH key IDs",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"tailscale_config_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Tailscale config ID",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Database status",
			},
			"database_hostname": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Hostname for connecting to the database",
			},
			"database_service_port": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Service port for connecting to the database",
			},
			"database_username": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Database username to use when connecting",
			},
			"database_password": {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "Database password to use when connecting",
			},
			"stripe_checkout_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Stripe checkout URL to complete billing when required",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Database creation timestamp",
			},
		},
	}
}

func resourceDatabaseCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)
	var diags diag.Diagnostics

	// Preparar el request body
	plan := d.Get("database_plan").([]interface{})[0].(map[string]interface{})
	billableItemsSet := plan["billable_items"].(*schema.Set)
	billableItems := billableItemsSet.List()

	databasePlanBI := make([]map[string]interface{}, len(billableItems))
	for i, item := range billableItems {
		itemMap := item.(map[string]interface{})
		databasePlanBI[i] = map[string]interface{}{
			"billableItemId": itemMap["billable_item_id"].(string),
			"quantity":       itemMap["quantity"].(int),
		}
	}

	requestBody := map[string]interface{}{
		"organizationSlug": d.Get("organization_slug").(string),
		"namespaceSlug":    d.Get("namespace_slug").(string),
		"engineId":         d.Get("engine_id").(string),
		"regionId":         d.Get("region_id").(string),
		"details": map[string]interface{}{
			"name":        d.Get("name").(string),
			"description": d.Get("description").(string),
		},
		"databasePlanDetails": map[string]interface{}{
			"databasePlanBI": databasePlanBI,
		},
	}

	if v, ok := d.GetOk("ip_whitelist_ids"); ok {
		requestBody["ipWhitelistIds"] = v.([]interface{})
	}

	if v, ok := d.GetOk("ssh_key_ids"); ok {
		requestBody["sshKeyIds"] = v.([]interface{})
	}

	if v, ok := d.GetOk("tailscale_config_id"); ok {
		requestBody["tailscaleConfigId"] = v.(string)
	}

	resp, err := c.Post("/api/v1/databases", requestBody)
	if err != nil {
		return diag.FromErr(err)
	}

	// Extraer el ID de la base de datos creada
	data := resp.Data.(map[string]interface{})
	database := data["database"].(map[string]interface{})
	// Manejar id que puede ser string o float64
	var databaseId string
	switch v := database["id"].(type) {
	case float64:
		databaseId = fmt.Sprintf("%.0f", v)
	case string:
		databaseId = v
	default:
		databaseId = fmt.Sprintf("%v", v)
	}

	d.SetId(databaseId)

	if url := extractStripeCheckoutURL(data); url != "" {
		// Imprimir directamente a /dev/tty para que sea visible sin TF_LOG
		if tty, err := os.OpenFile("/dev/tty", os.O_WRONLY, 0); err == nil {
			fmt.Fprintf(tty, "\n")
			fmt.Fprintf(tty, "╔════════════════════════════════════════════════════════════╗\n")
			fmt.Fprintf(tty, "║  ⚠️  PAYMENT REQUIRED                                      ║\n")
			fmt.Fprintf(tty, "╚════════════════════════════════════════════════════════════╝\n")
			fmt.Fprintf(tty, "\n")
			fmt.Fprintf(tty, "The database requires payment to continue provisioning.\n")
			fmt.Fprintf(tty, "Please open this URL to complete the Stripe checkout:\n\n")
			fmt.Fprintf(tty, "  %s\n\n", url)
			fmt.Fprintf(tty, "Waiting for payment completion...\n")
			fmt.Fprintf(tty, "\n")
			tty.Close()
		}

		tflog.Info(ctx, "Database provisioning blocked until Stripe checkout completes", map[string]interface{}{
			"stripe_checkout_url": url,
			"database_id":         databaseId,
		})
		if err := d.Set("stripe_checkout_url", url); err != nil {
			return diag.FromErr(err)
		}
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  "Payment required",
			Detail:   fmt.Sprintf("Open the checkout URL to complete billing and resume provisioning: %s", url),
		})
	}

	if _, err := waitForDatabaseCredentials(ctx, c, databaseId); err != nil {
		return diag.FromErr(err)
	}

	readDiags := resourceDatabaseRead(ctx, d, m)
	diags = append(diags, readDiags...)
	return diags
}

func resourceDatabaseRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)

	resp, err := c.Get("/api/v1/databases/" + d.Id())
	if err != nil {
		if apiErr, ok := err.(*client.APIError); ok && apiErr.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	data := resp.Data.(map[string]interface{})

	// Helper para convertir id a string
	idToString := func(v interface{}) string {
		switch val := v.(type) {
		case float64:
			return fmt.Sprintf("%.0f", val)
		case string:
			return val
		default:
			return fmt.Sprintf("%v", val)
		}
	}

	d.Set("name", data["name"])
	d.Set("description", data["description"])
	d.Set("status", data["status"])
	d.Set("engine_id", idToString(data["engineId"]))
	d.Set("region_id", idToString(data["regionId"]))

	if createdAt, ok := data["createdAt"]; ok {
		d.Set("created_at", createdAt)
	}

	if url := extractStripeCheckoutURL(data); url != "" {
		d.Set("stripe_checkout_url", url)
	} else {
		d.Set("stripe_checkout_url", "")
	}

	params := mapDatabaseParams(data["databaseParams"])
	d.Set("database_hostname", params["database_hostname"])
	d.Set("database_service_port", params["database_service_port"])

	username, password := selectDatabaseUser(data["databaseUsers"])
	d.Set("database_username", username)
	d.Set("database_password", password)

	return nil
}

func resourceDatabaseUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Por ahora, solo actualizamos los campos que se pueden modificar
	// La actualización completa requeriría un endpoint PUT/PATCH en la API
	return resourceDatabaseRead(ctx, d, m)
}

func resourceDatabaseDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)

	_, err := c.Delete("/api/v1/databases/" + d.Id())
	if err != nil {
		if apiErr, ok := err.(*client.APIError); ok && apiErr.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}

func waitForDatabaseCredentials(ctx context.Context, c *client.Client, databaseId string) (map[string]interface{}, error) {
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"creating", "deploying", "waiting_credentials", "billing_pending"},
		Target:     []string{"deployed"},
		MinTimeout: 5 * time.Second,
		Delay:      5 * time.Second,
		Timeout:    30 * time.Minute,
		Refresh: func() (interface{}, string, error) {
			resp, err := c.Get("/api/v1/databases/" + databaseId)
			if err != nil {
				return nil, "", err
			}

			data, ok := resp.Data.(map[string]interface{})
			if !ok {
				return nil, "", fmt.Errorf("unexpected database response format")
			}

			status, _ := data["status"].(string)
			if url := extractStripeCheckoutURL(data); url != "" {
				tflog.Warn(ctx, "Waiting for user to complete Stripe checkout", map[string]interface{}{
					"stripe_checkout_url": url,
					"database_id":         databaseId,
				})
			}
			if credentialsAreReady(data) {
				return data, status, nil
			}

			return data, "waiting_credentials", nil
		},
	}

	result, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return nil, err
	}

	data, ok := result.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected result type %T when waiting for database %s", result, databaseId)
	}

	return data, nil
}

func credentialsAreReady(data map[string]interface{}) bool {
	params := mapDatabaseParams(data["databaseParams"])
	if params["database_hostname"] == "" || params["database_service_port"] == "" {
		return false
	}

	username, password := selectDatabaseUser(data["databaseUsers"])
	return username != "" && password != ""
}

func mapDatabaseParams(raw interface{}) map[string]string {
	result := map[string]string{
		"database_hostname":     "",
		"database_service_port": "",
	}

	params, ok := raw.([]interface{})
	if !ok {
		return result
	}

	for _, p := range params {
		paramMap, ok := p.(map[string]interface{})
		if !ok {
			continue
		}

		key, _ := paramMap["key"].(string)
		value, _ := paramMap["value"].(string)

		if key == "" {
			continue
		}
		result[key] = value
	}

	return result
}

func selectDatabaseUser(raw interface{}) (string, string) {
	users, ok := raw.([]interface{})
	if !ok || len(users) == 0 {
		return "", ""
	}

	var fallbackUsername, fallbackPassword string
	for _, u := range users {
		userMap, ok := u.(map[string]interface{})
		if !ok {
			continue
		}

		username, _ := userMap["username"].(string)
		password, _ := userMap["password"].(string)
		role, _ := userMap["role"].(string)

		if username == "" || password == "" {
			continue
		}

		if role == "root" || username == "root" {
			return username, password
		}

		if fallbackUsername == "" {
			fallbackUsername = username
			fallbackPassword = password
		}
	}

	return fallbackUsername, fallbackPassword
}

func extractStripeCheckoutURL(data map[string]interface{}) string {
	if sessionRaw, ok := data["stripeCheckoutSession"]; ok && sessionRaw != nil {
		if sessionMap, ok := sessionRaw.(map[string]interface{}); ok {
			if url, _ := sessionMap["url"].(string); url != "" {
				return url
			}
		}
	}
	return ""
}
