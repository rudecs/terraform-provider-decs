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

	// "time"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

func diskSubresourceSchema() map[string]*schema.Schema {
	rets := map[string]*schema.Schema {
		"label": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Unique label to identify this disk among other disks connected to this VM.",
		},

		"size": {
			Type:        schema.TypeInt,
			Required:    true,
			ValidateFunc: validation.IntAtLeast(1),
			Description: "Size of the disk in GB.",
		},

		"pool": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "default",
			Description: "Pool from which this disk should be provisioned.",
		},

		"provider": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "default",
			Description: "Storage provider (storage technology type) by which this disk should be served.",
		},
		
	}

	return rets
}
