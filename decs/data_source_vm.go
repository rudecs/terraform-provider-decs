/*
Copyright (c) 2019-2021 Digital Energy Cloud Solutions LLC. All Rights Reserved.
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
	// "net/url"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

func flattenVm(d *schema.ResourceData, vm_facts string) error {
	// NOTE: this function modifies ResourceData argument - as such it should never be called
	// from resourceVmExists(...) method
	model := MachinesGetResp{}
	log.Printf("flattenVm: ready to unmarshal string %q", vm_facts) 
	err := json.Unmarshal([]byte(vm_facts), &model)
	if err != nil {
		return err
	}

	log.Printf("flattenVm: model.ID %d, model.ResGroupID %d", model.ID, model.ResGroupID)
			   
	d.SetId(fmt.Sprintf("%d", model.ID))
	d.Set("name", model.Name)
	d.Set("rgid", model.ResGroupID)
	d.Set("cpu", model.Cpu)
	d.Set("ram", model.Ram)
	// d.Set("boot_disk", model.BootDisk)
	d.Set("image_id", model.ImageID)
	d.Set("description", model.Description)

	bootdisk_map := make(map[string]interface{})
	bootdisk_map["size"] = model.BootDisk
	bootdisk_map["label"] = "boot"
	bootdisk_map["pool"] = "default"
	bootdisk_map["provider"] = "default"
	
	if err = d.Set("boot_disk", []interface{}{bootdisk_map}); err != nil {
		return err
	}

	if len(model.DataDisks) > 0 {
		log.Printf("flattenVm: calling flattenDataDisks")
		if err = d.Set("data_disks", flattenDataDisks(model.DataDisks)); err != nil {
			return err
		}
	}

	if len(model.NICs) > 0 {
		log.Printf("flattenVm: calling flattenNICs")
		if err = d.Set("nics", flattenNICs(model.NICs)); err != nil {
			return err
		}
		log.Printf("flattenVm: calling flattenNetworks")
		if err = d.Set("networks", flattenNetworks(model.NICs)); err != nil {
			return err
		}
	}

	if len(model.GuestLogins) > 0 {
		log.Printf("flattenVm: calling flattenGuestLogins")
		guest_logins := flattenGuestLogins(model.GuestLogins)
		if err = d.Set("guest_logins", guest_logins); err != nil {
			return err
		}
		// set user & password attributes to the corresponding values of the 1st item in the list
		if err = d.Set("user", guest_logins[0]["user"]); err != nil {
			return err
		}
		if err = d.Set("password", guest_logins[0]["password"]); err != nil {
			return err
		}
	}

	return nil
}

func dataSourceVmRead(d *schema.ResourceData, m interface{}) error {
	vm_facts, err := utilityVmCheckPresence(d, m)
	if vm_facts == "" {
		// if empty string is returned from utilityVmCheckPresence then there is no
		// such VM and err tells so - just return it to the calling party 
		d.SetId("") // ensure ID is empty
		return err
	}

	return flattenVm(d, vm_facts)
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

			/*
			"image_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the OS image this virtual machine is based on.",
			},
			*/

			"boot_disk": {
				Type:        schema.TypeList,
				Computed:    true,
				MinItems:    1,
				Elem:        &schema.Resource {
					Schema:  diskSubresourceSchema(),
				},
				Description: "Specification for a boot disk on this virtual machine.",
			},

			"data_disks": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Resource {
					Schema:  diskSubresourceSchema(),
				},
				Description: "Specification for data disks on this virtual machine.",
			},

			"guest_logins": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Resource {
					Schema:  loginsSubresourceSchema(),
				},
				Description: "Specification for guest logins on this virtual machine.",
			},

			"networks": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Resource {
					Schema:  networkSubresourceSchema(),
				},
				Description: "Specification for the networks to connect this virtual machine to.",
			},

			"nics": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Resource {
					Schema:  nicSubresourceSchema(),
				},
				Description: "Specification for the virutal NICs allocated to this virtual machine.",
			},

			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Description of this virtual machine.",
			},

			"user": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Default login name for the guest OS on this virtual machine.",
			},

			"password": {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "Default password for the guest OS login on this virtual machine.",
			},

		},
	}
}