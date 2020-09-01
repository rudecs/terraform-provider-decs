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
	"strconv"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)


func resourceVmCreate(d *schema.ResourceData, m interface{}) error {
	machine := &MachineConfig{
		ResGroupID:       d.Get("rgid").(int),
		Name:             d.Get("name").(string),
		Cpu:              d.Get("cpu").(int),
		Ram:              d.Get("ram").(int),
		ImageID:          d.Get("image_id").(int),
		Description:      d.Get("description").(string),
	}
	// BootDisk
	// DataDisks
	// Networks
	// PortForwards
	// SshKeyData string
	log.Printf("resourceVmCreate: called for VM name %q, ResGroupID %d", machine.Name, machine.ResGroupID)
	
	var subres_list []interface{}
	var subres_data map[string]interface{}
	var arg_value interface{}
	var arg_set bool
	// boot disk list is a required argument and has only one element,
	// which is of type diskSubresourceSchema
	subres_list = d.Get("boot_disk").([]interface{})
	subres_data = subres_list[0].(map[string]interface{})
	machine.BootDisk.Label = subres_data["label"].(string)
	machine.BootDisk.Size = subres_data["size"].(int)
	machine.BootDisk.Pool = subres_data["pool"].(string)
	machine.BootDisk.Provider = subres_data["provider"].(string)

	
	arg_value, arg_set = d.GetOk("data_disks")
	if arg_set {
		log.Printf("resourceVmCreate: calling makeDisksConfig")
		machine.DataDisks, _ = makeDisksConfig(arg_value.([]interface{}))
	}
	
	arg_value, arg_set = d.GetOk("networks")
	if arg_set {
		log.Printf("resourceVmCreate: calling makeNetworksConfig")
		machine.Networks, _ = makeNetworksConfig(arg_value.([]interface{}))
	}

	arg_value, arg_set = d.GetOk("port_forwards")
	if arg_set {
		log.Printf("resourceVmCreate: calling makePortforwardsConfig")
		machine.PortForwards, _ = makePortforwardsConfig(arg_value.([]interface{}))
	}
	
	arg_value, arg_set = d.GetOk("ssh_keys")
	if arg_set {
		log.Printf("resourceVmCreate: calling makeSshKeysConfig")
		machine.SshKeys, _ = makeSshKeysConfig(arg_value.([]interface{}))
	}
	
	// create basic VM (i.e. without port forwards and ext network connections - those will be done
	// by separate API calls)
	d.Partial(true)
	controller := m.(*ControllerCfg)
	url_values := &url.Values{}
	url_values.Add("cloudspaceId", fmt.Sprintf("%d", machine.ResGroupID))
	url_values.Add("name", machine.Name)
	url_values.Add("description", machine.Description)
	url_values.Add("vcpus", fmt.Sprintf("%d", machine.Cpu))
	url_values.Add("memory", fmt.Sprintf("%d", machine.Ram))
	url_values.Add("imageId", fmt.Sprintf("%d", machine.ImageID))
	url_values.Add("disksize", fmt.Sprintf("%d", machine.BootDisk.Size))
	if len(machine.SshKeys) > 0 {
		url_values.Add("userdata", makeSshKeysArgString(machine.SshKeys))
	}
	api_resp, err := controller.decsAPICall("POST", MachineCreateAPI, url_values)
	if err != nil {
		return err
	}
	d.SetId(api_resp) // machines/create API plainly returns ID of the new VM on success
	machine.ID, _ = strconv.Atoi(api_resp)
	d.SetPartial("name")
	d.SetPartial("description")
	d.SetPartial("cpu")
	d.SetPartial("ram")
	d.SetPartial("image_id")
	d.SetPartial("boot_disk")
	if len(machine.SshKeys) > 0 {
		d.SetPartial("ssh_keys")
	}

	log.Printf("resourceVmCreate: new VM ID %d, name %q created", machine.ID, machine.Name)

	if len(machine.DataDisks) > 0 || len(machine.PortForwards) > 0 {
		// for data disk or port foreards provisioning we have to know Tenant ID
		// and Grid ID so we call utilityResgroupConfigGet method to populate these 
		// fields in the machine structure that will be passed to provisionVmDisks or
		// provisionVmPortforwards
		log.Printf("resourceVmCreate: calling utilityResgroupConfigGet")
		resgroup, err := controller.utilityResgroupConfigGet(machine.ResGroupID)
		if err == nil {
			machine.TenantID = resgroup.TenantID
			machine.GridID = resgroup.GridID
			machine.ExtIP = resgroup.ExtIP
			log.Printf("resourceVmCreate: tenant ID %d, GridID %d, ExtIP %q", 
			machine.TenantID, machine.GridID, machine.ExtIP)
		}
	}

	//
	// Configure data disks
	disks_ok := true
	if len(machine.DataDisks) > 0 {
		log.Printf("resourceVmCreate: calling utilityVmDisksProvision for disk count %d", len(machine.DataDisks))
		if machine.TenantID == 0 {
			// if TenantID is still 0 it means that we failed to get Resgroup Facts by
			// a previous call to utilityResgroupGetFacts,
			// hence we do not have technical ability to provision data disks
			disks_ok = false
		} else {
			// provisionVmDisks accomplishes two steps for each data disk specification
			// 1) creates the disks
			// 2) attaches them to the VM
			err = controller.utilityVmDisksProvision(machine)
			if err != nil {
				disks_ok = false
			}
		}
	}

	if disks_ok {
		d.SetPartial("data_disks")
	}
	
	//
	// Configure port forward rules
	pfws_ok := true
	if len(machine.PortForwards) > 0 {
		log.Printf("resourceVmCreate: calling utilityVmPortforwardsProvision for pfw rules count %d", len(machine.PortForwards))
		if machine.ExtIP == "" {
			// if ExtIP is still empty it means that we failed to get Resgroup Facts by
			// a previous call to utilityResgroupGetFacts,
			// hence we do not have technical ability to provision port forwards
			pfws_ok = false
		} else {
			err := controller.utilityVmPortforwardsProvision(machine)
			if err != nil {
				pfws_ok = false
			}	
		}
	}
	if pfws_ok {
		//  there were no errors reported when configuring port forwards
		d.SetPartial("port_forwards")
	}

	//
	// Configure external networks
	// NOTE: currently only one external network can be attached to each VM, so in the current
	// implementation we ignore all but the 1st network definition
	nets_ok := true
	if len(machine.Networks) > 0 {
		log.Printf("resourceVmCreate: calling utilityVmNetworksProvision for networks count %d", len(machine.Networks))
		err := controller.utilityVmNetworksProvision(machine)
		if err != nil {
			nets_ok = false
		}
	}
	if nets_ok {
		// there were no errors reported when configuring networks
		d.SetPartial("networks")
	}

	if ( disks_ok && nets_ok && pfws_ok ) {
		// if there were no errors in setting any of the subresources, we may leave Partial mode
		d.Partial(false)
	}

	// resourceVmRead will also update resource ID on success, so that Terraform will know
	// that resource exists
	return resourceVmRead(d, m)
}

func resourceVmRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("resourceVmRead: called for VM name %q, ResGroupID %d", 
	           d.Get("name").(string), d.Get("rgid").(int))
	
	vm_facts, err := utilityVmCheckPresence(d, m)
	if vm_facts == "" {
		if err != nil {
			return err
		}
		// VM was not found
		return nil
	}

	if err = flattenVm(d, vm_facts); err != nil {
		return err
	}
	log.Printf("resourceVmRead: after flattenVm: VM ID %s, VM name %q, ResGroupID %d", 
	           d.Id(), d.Get("name").(string), d.Get("rgid").(int))

	// Not all parameters, that we may need, are returned by machines/get API
	// Continue with further reading of VM subresource parameters:
	controller := m.(*ControllerCfg)
	url_values := &url.Values{}

	/*
	// Obtain information on external networks
	url_values.Add("machineId", d.Id())
	body_string, err := controller.decsAPICall("POST", VmExtNetworksListAPI, url_values)
	if err != nil {
		return err
	}

	net_list := ExtNetworksResp{}
	err = json.Unmarshal([]byte(body_string), &net_list)
	if err != nil {
		return err
	}

	if len(net_list) > 0 {
		if err = d.Set("networks", flattenNetworks(net_list)); err != nil {
			return err
		}
	}
	*/

	/*
	// Ext networks flattening is now done inside flattenVm because it is currently based
	// on data read into NICs component by machine/get API call

	if err = d.Set("networks", flattenNetworks()); err != nil {
		return err
	}
	*/

	//
	// Obtain information on port forwards
	url_values.Add("cloudspaceId", fmt.Sprintf("%d",d.Get("rgid")))
	url_values.Add("machineId", d.Id())
	pfw_list := PortforwardsResp{}
	body_string, err := controller.decsAPICall("POST", PortforwardsListAPI, url_values)
	if err != nil {
		return err
	}
	err = json.Unmarshal([]byte(body_string), &pfw_list)
	if err != nil {
		return err
	}

	if len(pfw_list) > 0 {
		if err = d.Set("port_forwards", flattenPortforwards(pfw_list)); err != nil {
			return err
		}
	}

	return nil
}

func resourceVmUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("resourceVmUpdate: called for VM name %q, ResGroupID %d", 
			   d.Get("name").(string), d.Get("rgid").(int))
			   
	return resourceVmRead(d, m)
}

func resourceVmDelete(d *schema.ResourceData, m interface{}) error {
	// NOTE: this method destroys target VM with flag "permanently", so there is no way to
	// restore destroyed VM
	log.Printf("resourceVmDelete: called for VM name %q, ResGroupID %d", 
	           d.Get("name").(string), d.Get("rgid").(int))
			   
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
	log.Printf("resourceVmExist: called for VM name %q, ResGroupID %d", 
			   d.Get("name").(string), d.Get("rgid").(int))
			   
	vm_facts, err := utilityVmCheckPresence(d, m)
	if vm_facts == "" {
		if err != nil {
			return false, err
		}
		return false, nil
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
				Description:  "Name of this virtual machine. This parameter is case sensitive.",
			},

			"rgid": {
				Type:         schema.TypeInt,
				Required:     true,
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
				Optional:    true,
				MaxItems:    8,
				Elem:        &schema.Resource {
					Schema:  networkSubresourceSchema(),
				},
				Description: "Specification for the networks to connect this virtual machine to.",
			},

			"nics": {
				Type:        schema.TypeList,
				Computed:    true,
				MaxItems:    8,
				Elem:        &schema.Resource {
					Schema:  nicSubresourceSchema(),
				},
				Description: "Specification for the virutal NICs allocated to this virtual machine.",
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