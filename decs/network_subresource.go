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

func makeNetworksConfig(arg_list []interface{}) (nets []NetworkConfig, count int) {
	count = len(arg_list) 
	if count < 1 {
		return nil, 0
	}
	
	// allocate Networks list and fill it 
	nets = make([]NetworkConfig, count)
	var subres_data map[string]interface{}
	for index, value := range arg_list {
			subres_data = value.(map[string]interface{})
			nets[index].Label = subres_data["label"].(string)
			nets[index].NetworkID = subres_data["id"].(int)
		}

	return nets, count
} 

func flattenNetworks(nets []ExtNetworkRecord) []interface{} {
	result := make([]interface{}, len(nets))
	elem := make(map[string]interface{})

	for index, value := range nets {
		elem["network_id"] = value.ID
		elem["ip_range"] = value.IPRange
		result[index] = elem
	}

	return result 
}

func networkSubresourceSchema() map[string]*schema.Schema {
	rets := map[string]*schema.Schema {
		"network_id": {
			Type:        schema.TypeInt,
			Required:    true,
			ValidateFunc: validation.IntAtLeast(1),
			Description: "ID of the network to attach to this VM.",
		},

		"label": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Unique label of this network connection to identify it amnong other connections for this VM.",
		},

		"ip_range": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Range of IP addresses defined for this network.",
		},
	}

	return rets
}

func makePortforwardsConfig(arg_list []interface{}) (pfws []PortforwardConfig, count int) {
	count = len(arg_list) 
	if count < 1 {
		return nil, 0
	}

	pfws = make([]PortforwardConfig, count)
	var subres_data map[string]interface{}
	for index, value := range arg_list {
		subres_data = value.(map[string]interface{})
		pfws[index].Label = subres_data["label"].(string)
		pfws[index].ExtPort = subres_data["ext_port"].(int)
		pfws[index].IntPort = subres_data["int_port"].(int)
		pfws[index].Proto = subres_data["proto"].(string)
	}

	return pfws, count
}

func flattenPortforwards(pfws []PortforwardRecord) []interface{} {
	result := make([]interface{}, len(pfws))
	elem := make(map[string]interface{})

	for index, value := range pfws {
		elem["label"] = "default"
		elem["ext_port"] = value.ExtPort
		elem["int_port"] = value.IntPort
		elem["proto"] = value.Proto
		elem["ext_ip"] = value.ExtIP
		elem["int_ip"] = value.IntIP
		result[index] = elem
	}

	return result 
}

func portforwardSubresourceSchema() map[string]*schema.Schema {
	rets := map[string]*schema.Schema {
		"label": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Unique label of this network connection to identify it amnong other connections for this VM.",
		},
		
		"ext_port": {
			Type:        schema.TypeInt,
			Required:    true,
			ValidateFunc: validation.IntBetween(1, 65535),
			Description: "External port number for this port forwarding rule.",
		},

		"int_port": {
			Type:        schema.TypeInt,
			Required:    true,
			ValidateFunc: validation.IntBetween(1, 65535),
			Description: "Internal port number for this port forwarding rule.",
		},

		"proto": {
			Type:        schema.TypeInt,
			Required:    true,
			// ValidateFunc: validation.IntBetween(1, ),
			Description: "Protocol type for this port forwarding rule. Should be either 'tcp' or 'udp'.",
		},

		"ext_ip": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: ".",
		},

		"int_ip": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: ".",
		},
	}

	return rets
}

func flattenNICs(nics []NicRecord) []interface{} {
	var result = make([]interface{}, len(nics))
	elem := make(map[string]interface{})

	for index, value := range nics {
		elem["status"] = value.Status
		elem["type"] = value.NicType
		elem["mac"] = value.MacAddress
		elem["ip_address"] = value.IPAddress
		elem["parameters"] = value.Params
		elem["reference_id"] = value.ReferenceID
		elem["network_id"] = value.NetworkID
		result[index] = elem
	}

	return result
}

func nicSubresourceSchema() map[string]*schema.Schema {
	rets := map[string]*schema.Schema {
		"status": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Current status of this NIC.",
		},

		"type": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Type of this NIC.",
		},
		
		"mac": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "MAC address assigned to this NIC.",
		},

		"ip_address": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "IP address assigned to this NIC.",
		},

		"parameters": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Additional NIC parameters.",
		},

		"reference_id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Reference ID of this NIC.",
		},

		"network_id": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Network ID which this NIC is connected to.",
		},
	}

	return rets
}
