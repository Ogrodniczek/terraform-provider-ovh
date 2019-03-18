package ovh

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
)

type OvhKubeNode struct {
	CreatedAt  string `json:"createdAt"`
	Id         string `json:"id"`
	InstanceId string `json:"instanceId,omitempty"`
	Name       string `json:"name,omitempty"`
	Version    string `json:"version"`
	IsUpToDate bool   `json:"isUpToDate"`
	UpdatedAt  string `json:"updatedAt"`
	Flavor     string `json:"flavor"`
	Status     string `json:"status"`
	ProjectId  string `json:"projectId"`
}

func resourceOvhKubeNode() *schema.Resource {
	return &schema.Resource{
		Create: resourceOvhDomainZoneRedirectionCreate,
		Read:   resourceOvhDomainZoneRedirectionRead,
		Delete: resourceOvhDomainZoneRedirectionDelete,

		Schema: map[string]*schema.Schema{
			"createdAt": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"instanceId": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"version": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"isUpToDate": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"updatedAt": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"flavor": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					err := validateStringEnum(v.(string), []string{
						"INSTALLING", "UPDATING", "RESETTING", "SUSPENDING", "REOPENING", "DELETING", "SUSPENDED",
						"ERROR", "USER_ERROR", "USER_QUOTA_ERROR", "USER_NODE_NOT_FOUND_ERROR", "USER_NODE_SUSPENDED_SERVICE", "READY",
					})
					if err != nil {
						errors = append(errors, err)
					}
					return
				},
			},
			"projectId": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceOvhKubeNodeCreate(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*Config)

	newKubeNode := &OvhKubeNode{
		ProjectId: d.Get("projectId").(string),
		Flavor:    d.Get("flavor").(string),
		Name:      d.Get("name").(string),
	}

	log.Printf("[DEBUG] OVH Kubernetes node create configuration: %#v", newKubeNode)

	resultKubeNode := OvhKubeNode{}

	err := provider.OVHClient.Post(
		fmt.Sprintf("/kube/%s/publiccloud/node", d.Get("projectId").(string)),
		newKubeNode,
		&resultKubeNode,
	)

	if err != nil {
		return fmt.Errorf("Failed to create OVH Kubernetes Node: %s", err)
	}

	d.SetId(resultKubeNode.Id)

	log.Printf("[INFO] OVH Kubernetes Node ID: %s", d.Id())

	//	if err := ovhDomainZoneRefresh(d, meta); err != nil {
	//	log.Printf("[WARN] OVH Domain zone refresh after redirection creation failed: %s", err)
	//}

	return resourceOvhKubeNodeRead(d, meta)
}

func resourceOvhKubeNodeRead(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*Config)

	kubenode := OvhKubeNode{}
	err := provider.OVHClient.Get(
		fmt.Sprintf("/kube/%s/publiccloud/node/%s", d.Get("projectId").(string), d.Id()),
		&kubenode,
	)

	if err != nil {
		d.SetId("")
		return nil
	}

	d.Set("createdAt", kubenode.CreatedAt)
	d.Set("id", kubenode.Id)
	d.Set("instanceId", kubenode.InstanceId)
	d.Set("name", kubenode.Name)
	d.Set("version", kubenode.Version)
	d.Set("isUpToDate", kubenode.IsUpToDate)
	d.Set("updatedAt", kubenode.UpdatedAt)
	d.Set("flavor", kubenode.Flavor)
	d.Set("status", kubenode.Status)
	d.Set("projectId", kubenode.ProjectId)

	return nil
}

func resourceOvhKubeNodeDelete(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*Config)

	log.Printf("[INFO] Deleting OVH Kubernetes node in: %s, id: %s", d.Get("serviceName").(string), d.Id())

	err := provider.OVHClient.Delete(
		fmt.Sprintf("/kube/%s/publiccloud/node/%s", d.Get("projectId").(string), d.Id()),
		nil,
	)

	if err != nil {
		return fmt.Errorf("Error deleting OVH Kubernetes node: %s", err)
	}

	//if err := ovhDomainZoneRefresh(d, meta); err != nil {
	//	log.Printf("[WARN] OVH Domain zone refresh after redirection deletion failed: %s", err)
	//}

	return nil
}
