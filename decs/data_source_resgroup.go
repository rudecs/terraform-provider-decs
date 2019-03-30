/*
Copyright (c) 2019 Digital Energy Cloud Solutions LLC. All Rights Reserved.
Author: Sergey Shubin, <sergey.shubin@digitalenergy.online>, <svs1370@gmail.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package decs

import (

	"encoding/json"
	"fmt"
	"log"
	"net/url"

	"github.com/hashicorp/terraform/helper/schema"
	// "github.com/hashicorp/terraform/helper/validation"

)


func dataSourceResgroupRead(d *schema.ResourceData, m interface{}) error {
	name := d.Get("name").(string)
	tenant_name := d.Get("tenant").(string)

	controller := m.(*ControllerCfg)
	url_values := &url.Values{}
	body_string, err := controller.decsAPICall("POST", CloudspacesListAPI, url_values)
	if err != nil {
		return err
	}

	log.Printf("%s", body_string)
	log.Printf("dataSourceResgroupRead: ready to decode response body from %q", CloudspacesListAPI)
	model := CloudspacesListResp{}
	err = json.Unmarshal([]byte(body_string), &model)
	if err != nil {
		return err
	}

	log.Printf("dataSourceResgroupRead: traversing decoded Json of length %d", len(model))
	for index, item := range model {
		// need to match VDC by name & tenant name
		if item.Name == name && item.TenantName == tenant_name {
			log.Printf("dataSourceResgroupRead: match ResGroup name %q / ID %d, tenant %q at index %d", 
			           item.Name, item.ID, item.TenantName, index)
			d.SetId(fmt.Sprintf("%d", item.ID))
			d.Set("name", item.Name)
			d.Set("tenant_id", item.TenantID)
			d.Set("grid_id", item.GridID)
			d.Set("public_ip", item.PublicIP)

			// not all required information is returned by cloudspaces/list API, so we need to initiate one more
			// call to cloudspaces/get to obtain extra data to complete Resource population.
			// Namely, we need to extract resource quota settings
			req_values := &url.Values{} 
			req_values.Add("cloudspaceId", fmt.Sprintf("%d", item.ID))
			body_string, err := controller.decsAPICall("POST", CloudspacesGetAPI, req_values)
			if err != nil {
				return err
			}
			log.Printf("%s", body_string)
			log.Printf("dataSourceResgroupRead: ready to decode response body from %q", CloudspacesGetAPI)
			details := CloudspacesGetResp{}
			err = json.Unmarshal([]byte(body_string), &details)
			if err != nil {
				return err
			}
			log.Printf("dataSourceResgroupRead: calling flattenQuotas()")
			if err = d.Set("quotas", flattenQuotas(details.Quotas)); err != nil {
				return err
			}

			return nil
		}
	}

	return fmt.Errorf("Cannot find resource group name %q owned by tenant %q", name, tenant_name)
}


func dataSourceResgroup() *schema.Resource {
	return &schema.Resource {
		SchemaVersion: 1,

		Read:   dataSourceResgroupRead,

		Timeouts: &schema.ResourceTimeout {
			Read:    &Timeout30s,
			Default: &Timeout60s,
		},

		Schema: map[string]*schema.Schema {
			"name": {
				Type:          schema.TypeString,
				Required:      true,
				Description:  "Name of this resource group. Names are case sensitive and unique within the context of a tenant.",
			},

			"tenant": &schema.Schema {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the tenant, which this resource group belongs to.",
			},

			"tenant_id": &schema.Schema {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Unique ID of the tenant, which this resource group belongs to.",
			},

			"grid_id": &schema.Schema {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Unique ID of the grid, where this resource group is deployed.",
			},

			"location": {
				Type:          schema.TypeString,
				Computed:      true,
				Description:  "Location of this resource group.",
			},

			"public_ip": {  // this may be obsoleted as new network segments and true resource groups are implemented
				Type:          schema.TypeString,
				Computed:      true,
				Description:  "Public IP address of this resource group (if any).",
			},

			"quotas": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Elem:        &schema.Resource {
					Schema:  quotasSubresourceSchema(),
				},
				Description: "Quotas on the resources for this resource group.",
			},
		},
	}
}

func flattenQuotas(quotas QuotaRecord) []interface{} {
	quotas_map :=  make(map[string]interface{})

	quotas_map["cpu"] = quotas.Cpu
	quotas_map["ram"] = int(quotas.Ram)
	quotas_map["disk"] = quotas.Disk
	quotas_map["ext_ips"] = quotas.ExtIPs

	result := make([]interface{}, 1)
	result[0] = quotas_map

	return result
}

func quotasSubresourceSchema() map[string]*schema.Schema {
	rets := map[string]*schema.Schema {
		"cpu": &schema.Schema {
			Type:        schema.TypeInt,
			Optional:    true,
			Default:     -1,
			Description: "The quota on the total number of CPUs in this resource group.",
		},

		"ram": &schema.Schema {
			Type:        schema.TypeInt, // NB: API expects this as float! This may be changed in the future.
			Optional:    true,
			Default:     -1,
			Description: "The quota on the total amount of RAM in this resource group, specified in MB.",
			},

		"disk": &schema.Schema {
			Type:        schema.TypeInt,
			Optional:    true,
			Default:     -1,
			Description: "The quota on the total volume of storage resources in this resource group, specified in GB.",
		},

		"ext_ips": &schema.Schema {
			Type:        schema.TypeInt,
			Optional:    true,
			Default:     -1,
			Description: "The quota on the total number of external IP addresses this resource group can use.",
		},
	}
	return rets
}