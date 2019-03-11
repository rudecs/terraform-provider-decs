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

func networkSubresourceSchema() map[string]*schema.Schema {
	rets := map[string]*schema.Schema {
		"id": {
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
	}

	return rets
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
	}

	return rets
}