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

	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

func flattenGuestLogins(logins []GuestLoginRecord) []interface{} {
	var length = 0

	result := make([]interface{}, length)
	if length == 0 {
		return result
	}

	elem := make(map[string]interface{})

	var subindex = 0
	for _, value := range logins {
		elem["guid"] = value.Guid
		elem["login"] = value.Login
		elem["password"] = value.Password
		result[subindex] = elem
		subindex += 1
	}

	return result
}

func loginsSubresourceSchema() map[string]*schema.Schema {
	rets := map[string]*schema.Schema {
		"guid": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "",
			Description: "GUID of this guest user.",
		},

		"login": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Login name of this guest user.",
		},

		"password": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Password of this guest user.",
		},
	}

	return rets
}
