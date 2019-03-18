package ovh

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceKube() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKubeRead,
		Schema: map[string]*schema.Schema{
			"url": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					err := validateStringEnum(v.(string), []string{
						"INSTALLING", "UPDATING", "RESETTING", "SUSPENDING", "REOPENING", "DELETING",
						"SUSPENDED", "ERROR", "USER_ERROR", "USER_QUOTA_ERROR", "READY",
					})
					if err != nil {
						errors = append(errors, err)
					}
					return
				},
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"nodes_url": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"update_policy": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"version": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"is_up_to_date": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
		},
	}
}

type Kube struct {
	Url                    string `json:"url"`
	Status                 string `json:"status"`
	Name                   string `json:"name"`
	NodesUrl               string `json:"nodesUrl"`
	CreatedAt              string `json:"createdAt"`
	UpdatePolicy           string `json:"updatePolicy"`
	Version                string `json:"version"`
	UpdatedAt              string `json:"updatedAt"`
	Id                     string `json:"id"`
	IsUpToDate             bool   `json:"isUpToDate"`
	ControlPlaneIsUpToDate bool   `json:"controlPlaneIsUpToDate"`
}

func dataSourceKubeRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	log.Printf("[DEBUG] Will list available Kubernetes clusters")

	response := []string{}
	err := config.OVHClient.Get("/kube", &response)

	if err != nil {
		return fmt.Errorf("Error calling /kube:\n\t %q", err)
	}

	kube := &Kube{}
	for _, serviceName := range response {
		err := config.OVHClient.Get(fmt.Sprintf("/kube/%s", serviceName), &kube)

		if err != nil {
			return fmt.Errorf("Error calling /kube/%s:\n\t %q", serviceName, err)
		}

		if v, ok := d.GetOk("url"); ok && v.(string) != kube.Url {
			continue
		}
		if v, ok := d.GetOk("status"); ok && v.(string) != kube.Status {
			continue
		}
		if v, ok := d.GetOk("name"); ok && v.(string) != kube.Name {
			continue
		}
		if v, ok := d.GetOk("nodes_url"); ok && v.(string) != kube.NodesUrl {
			continue
		}
		if v, ok := d.GetOk("created_at"); ok && v.(string) != kube.CreatedAt {
			continue
		}
		if v, ok := d.GetOk("update_policy"); ok && v.(string) != kube.UpdatePolicy {
			continue
		}
		if v, ok := d.GetOk("version"); ok && v.(string) != kube.Version {
			continue
		}
		if v, ok := d.GetOk("updated_at"); ok && v.(string) != kube.UpdatedAt {
			continue
		}
		if v, ok := d.GetOk("is_up_to_date"); ok && v.(bool) != kube.IsUpToDate {
			continue
		}
	}

	d.SetId(kube.Id)
	d.Set("url", kube.Url)
	d.Set("status", kube.Status)
	d.Set("name", kube.Name)
	d.Set("nodes_url", kube.NodesUrl)
	d.Set("created_at", kube.CreatedAt)
	d.Set("update_policy", kube.UpdatePolicy)
	d.Set("version", kube.Version)
	d.Set("updated_at", kube.UpdatedAt)
	d.Set("is_up_to_date", kube.IsUpToDate)

	return nil
}
