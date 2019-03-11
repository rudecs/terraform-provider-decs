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
	// "github.com/hashicorp/terraform/helper/validation"
)

func sshSubresourceSchema() map[string]*schema.Schema {
	rets := map[string]*schema.Schema {
		"user": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Name of the user on the guest OS of the new VM, for which the following SSH key will be authorized.",
		},

		"public_key": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Public part of SSH key to authorize to the specified user on the VM being created.",
		},
	}

	return rets
}

