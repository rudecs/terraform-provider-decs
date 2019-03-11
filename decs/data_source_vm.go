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
	// "log"
	// "net/url"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)


func dataSourceVmRead(d *schema.ResourceData, m interface{}) error {
	vm_facts, err := utilityVmCheckPresence(d, m)
	if vm_facts == "" {
		// if empty string is returned from utilityVmCheckPresence then there is no
		// such VM and err tells so - just return it to the calling party 
		return err
	}

	model := MachinesGetResp{}
	err = json.Unmarshal([]byte(vm_facts), &model)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%d", model.ID))
	d.Set("cpu", model.Cpu)
	d.Set("ram", model.Ram)
	d.Set("boot_disk", model.BootDisk)
	d.Set("image_id", model.ImageID)
	// d.Set("image_name", model.ImageName)
	d.Set("description", model.Description)
	// d.Set("field_name", value)
	return nil
}

func dataSourceVm() *schema.Resource {
	return &schema.Resource {
		SchemaVersion: 1,

		Read:   dataSourceVmRead,

		Timeouts: &schema.ResourceTimeout {
			Read:    &Timeout30s,
			Default: &Timeout60s,
		},

		Schema: map[string]*schema.Schema {
			"name": {
				Type:          schema.TypeString,
				Required:      true,
				Description:  "Name of this virtual machine. This parameter is case sensitive.",
			},

			"rgid": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntAtLeast(1),
				Description:  "ID of the resource group where this virtual machine is located.",
			},

			/*
			"internal_ip": {
				Type:          schema.TypeString,
				Computed:      true,
				Description:  "Internal IP address of this VM.",
			},
			*/

			"cpu": {
				Type:         schema.TypeInt,
				Computed:     true,
				Description:  "Number of CPUs allocated for this virtual machine.",
			},

			"ram": {
				Type:         schema.TypeInt,
				Computed:     true,
				Description:  "Amount of RAM in MB allocated for this virtual machine.",
			},

			"image_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "ID of the OS image this virtual machine is based on.",
			},

			"image_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the OS image this virtual machine is based on.",
			},

			"boot_disk": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Size of the boot disk on this virtual machine.",
			},

			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Description of this virtual machine.",
			},

		},
	}
}