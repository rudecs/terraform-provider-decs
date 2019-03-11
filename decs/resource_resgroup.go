/*
Copyright (c) 2019 Digital Energy Cloud Solutions. All Rights Reserved.

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
        "github.com/hashicorp/terraform/helper/schema"
)

func resourceResgroupCreate(d *schema.ResourceData, m interface{}) error {
	return resourceResgroupRead(d, m)
}

func resourceResgroupRead(d *schema.ResourceData, m interface{}) error {
	return nil // calling dataSourceResgroupRead(d, m) from here may not be the best idea - consider!
}

func resourceResgroupUpdate(d *schema.ResourceData, m interface{}) error {
	return resourceResgroupRead(d, m)
}

func resourceResgroupDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceResgroupExists(d *schema.ResourceData, m interface{}) (bool, error) {
	// Reminder: according to Terraform rules, this function should not modify ResourceData argument
	return true, nil
}

func resourceResgroup() *schema.Resource {
	return &schema.Resource {
		SchemaVersion: 1,

		Create: resourceResgroupCreate,
		Read:   resourceResgroupRead,
		Update: resourceResgroupUpdate,
		Delete: resourceResgroupDelete,
		Exists: resourceResgroupExists,

		Timeouts: &schema.ResourceTimeout {
			Create:  &Timeout180s,
			Read:    &Timeout30s,
			Update:  &Timeout180s,
			Delete:  &Timeout60s,
			Default: &Timeout60s,
		},

		Schema: map[string]*schema.Schema {
			"name": &schema.Schema {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of this resource group. Names are unique within the context of a tenant and case sensitive.",
			},

			"tenant": &schema.Schema {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the tenant, which this resource group belongs to.",
			},

			"location": &schema.Schema {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the location where this resource group should exist.",
			},

			"quota_cpu": &schema.Schema {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The quota on the total number of CPUs in this resource group.",
			},

			"quota_ram": &schema.Schema {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The quota on the total amount of RAM in this resource group, specified in MB.",
			},

			"quota_disk": &schema.Schema {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The quota on the total volume of storage resources in this resource group, specified in GB.",
			},
		},
	}
}