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
		"github.com/hashicorp/terraform/helper/schema"
		"github.com/hashicorp/terraform/helper/validation"
)

func resourceVmCreate(d *schema.ResourceData, m interface{}) error {
	return resourceVmRead(d, m)
}

func resourceVmRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceVmUpdate(d *schema.ResourceData, m interface{}) error {
	return resourceVmRead(d, m)
}

func resourceVmDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceVm() *schema.Resource {
	return &schema.Resource {
		Create: resourceVmCreate,
		Read:   resourceVmRead,
		Update: resourceVmUpdate,
		Delete: resourceVmDelete,

		Schema: map[string]*schema.Schema {
			"name": {
				Type:          schema.TypeString,
				Required:      true,
				ForceNew:      true,
				Description:  "Name of this virtual machine. This parameter is case sensitive.",
			},

			"rgid": {
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntAtLeast(1),
				Description:  "ID of the resource group where this virtual machine should be deployed.",
			},

			"cpu": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(1, 64),
				Description:  "Number of CPUs to allocate to this virtual machine.",
			},

			"ram": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntAtLeast(512),
				Description:  "Amount of RAM in MB to allocate to this virtual machine.",
			},

			"image_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the OS image to base this virtual machine on. This parameter is case sensitive.",
			},

			"boot_disk": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Elem:        &schema.Resource {
					Schema:  ??,
				},
				Description: "Specification for a boot disk on this virtual machine."
			},

			"networks": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    12,
				Elem:        &schema.Resource {
					Schema:  ??,
				},
				Description: "Specification for the networks to connect this virtual machine to."
			},
			
			"ssh_key": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    12,
				Elem:        &schema.Resource {
					Schema:  ??,
				},
				Description: "SSH keys to authorize on this virtual machine."
			},

			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of this virtual machine."
			},

			"id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Unique ID of this virtual machine. This parameter is assigned by the cloud when the machine is created."
			},
		},
	}
}