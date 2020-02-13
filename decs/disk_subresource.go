/*
Copyright (c) 2019-2020 Digital Energy Cloud Solutions LLC. All Rights Reserved.
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

	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

func makeDisksConfig(arg_list []interface{}) (disks []DiskConfig, count int) {
	// This method takes a list of disk definitions, coming from the corresponding
	// schema, and populates a list of DiskConfig structures according to these 
	// definitions.
	count = len(arg_list) 
	if count < 1 { 
		return nil, 0
	}

	// allocate DataDisks list and fill it 
	disks = make([]DiskConfig, count)
	var subres_data map[string]interface{}
	for index, value := range arg_list {
		subres_data = value.(map[string]interface{})
		disks[index].Label = subres_data["label"].(string)
		disks[index].Size = subres_data["size"].(int)
		disks[index].Pool = subres_data["pool"].(string)
		disks[index].SepId = subres_data["sep_id"].(int)
		disks[index].Description = subres_data["descr"].(string)
		disks[index].Type = subres_data["type"].(string)
	}

	return disks, count
}

func flattenDataDisks(disks []GenericDiskRecord) []interface{} {
	var length = 0
	for _, value := range disks {
		if value.Type == "D" {
			length += 1
		}
	}
	log.Printf("flattenDataDisks: found %d disks with D type", length)

	result := make([]interface{}, length)
	if length == 0 {
		return result
	}

	elem := make(map[string]interface{})

	var subindex = 0
	for _, value := range disks {
		if value.Type == "D" {
			elem["label"] = value.Label
			elem["size"] = value.SizeMax
			elem["disk_id"] = value.ID
			elem["pool"] = value.Pool
			elem["sep_id"] = value.SepId
			elem["type"] = value.Type
			elem["descr"] = value.Description
			result[subindex] = elem
			subindex += 1
		}
		
	}

	return result
}

/*
func makeDataDisksArgString(disks []DiskConfig) string {
	// Prepare a string with the sizes of data disks for the virtual machine.
	// It is designed to be passed as "datadisks" argument of virtual machine create API call.
	if len(disks) < 1 {
		return ""
	}
	return ""
}
*/

func diskSubresourceSchema() map[string]*schema.Schema {
	rets := map[string]*schema.Schema {
		"descr": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "",
			Description: "Description for this disk.",
		},

		"disk_id": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "ID of this disk resource.",
		},

		"label": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Unique label to identify this disk among other disks connected to this VM.",
		},

		"pool": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "default", // each SEP should implement "default" pool functionality
			Description: "Pool from which this disk should be provisioned.",
		},

		"sep_id": {
			Type:        schema.TypeInt,
			Optional:    true,
			Default:     1,
			Description: "ID of the Storage End-point provider (SEP) by which this disk should be served.",
		},

		"size": {
			Type:        schema.TypeInt,
			Required:    true,
			ValidateFunc: validation.IntAtLeast(1),
			Description: "Size of the disk in GB. For boot disks is should be enough to accomodate selected OS image.",
		},

		"type": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "D", // by default Data disks are created
			Description: "Type of disk to create. D for data (default), B for boot.",
		},
		
	}

	return rets
}
