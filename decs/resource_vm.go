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
	"net/url"
	// "strconv"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)


func resourceVmCreate(d *schema.ResourceData, m interface{}) error {
	machine := &MachineConfig{}
	machine.ResGroupID = d.Get("rgid").(int)
	machine.Name = d.Get("name").(string)
	machine.Cpu = d.Get("cpu").(int)
	machine.Ram = d.Get("ram").(int)
	machine.ImageID = d.Get("image_id").(int)
	// BootDisk
	// DataDisks
	// Networks
	// PortForwards
	// SshKeyData string
	machine.Description = d.Get("description").(string)
	
	var subres_list []interface{}
	var subres_data map[string]interface{}
	var get_value interface{}
	var arg_set bool
	// boot disk list is mandatory and has only one element, which is of type diskSubresourceSchema
	subres_list = d.Get("boot_disk").([]interface{})
	subres_data = subres_list[0].(map[string]interface{})
	machine.BootDisk.Label = subres_data["label"].(string)
	machine.BootDisk.Size = subres_data["size"].(int)
	machine.BootDisk.Pool = subres_data["pool"].(string)
	machine.BootDisk.Provider = subres_data["provider"].(string)

	
	get_value, arg_set = d.GetOk("data_disk")
	if arg_set && len(get_value.([]interface{})) > 0 {
		// allocate DataDisks list and fill it 
		machine.DataDisks = make([]DiskConfig, len(get_value.([]interface{})))
		for index, value := range get_value.([]interface{}) {
			subres_data = value.(map[string]interface{})
			machine.DataDisks[index].Label = subres_data["label"].(string)
			machine.DataDisks[index].Size = subres_data["size"].(int)
			machine.DataDisks[index].Pool = subres_data["pool"].(string)
			machine.DataDisks[index].Provider = subres_data["provider"].(string)
		}
	}
	
	get_value, arg_set = d.GetOk("networks")
	if arg_set && len(get_value.([]interface{})) > 0 {
		// allocate Networks list and fill it 
		machine.Networks = make([]NetworkConfig, len(get_value.([]interface{})))
		for index, value := range subres_list {
			subres_data = value.(map[string]interface{})
			machine.Networks[index].Label = subres_data["label"].(string)
			machine.Networks[index].NetworkID = subres_data["id"].(int)
		}
	}

	get_value, arg_set = d.GetOk("port_forwards")
	if arg_set && len(get_value.([]interface{})) > 0 {
		// allocate Portforwards list and fill it 
		machine.PortForwards = make([]PortforwardConfig, len(get_value.([]interface{})))
		for index, value := range subres_list {
			subres_data = value.(map[string]interface{})
			machine.PortForwards[index].Label = subres_data["label"].(string)
			machine.PortForwards[index].ExtPort = subres_data["ext_port"].(int)
			machine.PortForwards[index].IntPort = subres_data["int_port"].(int)
			machine.PortForwards[index].Proto = subres_data["proto"].(string)
		}
	}

	get_value, arg_set = d.GetOk("ssh_keys")
	if arg_set && len(get_value.([]interface{})) > 0 {
		// allocate Ssh keys list and fill it 
		machine.SshKeys = make([]SshKeyConfig, len(get_value.([]interface{})))
		for index, value := range subres_list {
			subres_data = value.(map[string]interface{})
			machine.SshKeys[index].User = subres_data["user"].(string)
			machine.SshKeys[index].SshKey = subres_data["public_key"].(string)
		}
	}
	
	controller := m.(*ControllerCfg)
	list_url_values := &url.Values{}
	list_url_values.Add("cloudspaceId", fmt.Sprintf("%d", machine.ResGroupID))
	list_url_values.Add("name", machine.Name)
	list_url_values.Add("description", machine.Description)
	list_url_values.Add("vcpus", fmt.Sprintf("%d", machine.Cpu))
	list_url_values.Add("memory", fmt.Sprintf("%d", machine.Ram))
	list_url_values.Add("imageId", fmt.Sprintf("%d", machine.ImageID))
	list_url_values.Add("disksize", fmt.Sprintf("%d", machine.BootDisk.Size))
	// list_url_values.Add("datadisks", build a list of data disk sizes)
	// list_url_values.Add("userdata", build a string for cloud init to deploy SSH keys to the new VM)
	_, err := controller.decsAPICall("POST", MachineCreateAPI, list_url_values)
	if err != nil {
		return err
	}

	return resourceVmRead(d, m)
}

func resourceVmRead(d *schema.ResourceData, m interface{}) error {
	vm_facts, err := utilityVmCheckPresence(d, m)
	if vm_facts == "" {
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
	// d.Set("boot_disk", model.BootDisk)
	d.Set("image_id", model.ImageID)
	d.Set("description", model.Description)

	return nil
}

func resourceVmUpdate(d *schema.ResourceData, m interface{}) error {
	return resourceVmRead(d, m)
}

func resourceVmDelete(d *schema.ResourceData, m interface{}) error {
	// NOTE: this method destroys target VM with flag "permanently", so there is no way to
	// restore destroyed VM
	vm_facts, err := utilityVmCheckPresence(d, m)
	if vm_facts == "" {
		// the target VM does not exist - in this case according to Terraform best practice 
		// we exit from Destroy method without error
		return nil
	}

	params := &url.Values{}
	params.Add("machineId", d.Id())
	params.Add("permanently", "true")

	controller := m.(*ControllerCfg)
	vm_facts, err = controller.decsAPICall("POST", MachineDeleteAPI, params)
	if err != nil {
		return err
	}

	return nil
}

func resourceVmExists(d *schema.ResourceData, m interface{}) (bool, error) {
	// Reminder: according to Terraform rules, this function should not modify its ResourceData argument
	vm_facts, err := utilityVmCheckPresence(d, m)
	if vm_facts == "" {
		return false, err
	}
	return true, nil
}

func resourceVm() *schema.Resource {
	return &schema.Resource {
		SchemaVersion: 1,

		Create: resourceVmCreate,
		Read:   resourceVmRead,
		Update: resourceVmUpdate,
		Delete: resourceVmDelete,
		Exists:  resourceVmExists,

		Timeouts: &schema.ResourceTimeout {
			Create:  &Timeout180s,
			Read:    &Timeout30s,
			Update:  &Timeout180s,
			Delete:  &Timeout60s,
			Default: &Timeout60s,
		},

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

			"image_id": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the OS image to base this virtual machine on.",
			},

			"boot_disk": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Elem:        &schema.Resource {
					Schema:  diskSubresourceSchema(),
				},
				Description: "Specification for a boot disk on this virtual machine.",
			},

			"data_disks": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    12,
				Elem:        &schema.Resource {
					Schema:  diskSubresourceSchema(),
				},
				Description: "Specification for data disks on this virtual machine.",
			},

			
			"networks": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    8,
				Elem:        &schema.Resource {
					Schema:  networkSubresourceSchema(),
				},
				Description: "Specification for the networks to connect this virtual machine to.",
			},
			
			"ssh_keys": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    12,
				Elem:        &schema.Resource {
					Schema:  sshSubresourceSchema(),
				},
				Description: "SSH keys to authorize on this virtual machine.",
			},

			"port_forwards": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    12,
				Elem:        &schema.Resource {
					Schema:  portforwardSubresourceSchema(),
				},
				Description: "Specification for the port forwards to configure for this virtual machine.",
			},

			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of this virtual machine.",
			},

		},
	}
}