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

	"json"
	"net/http"
	"net/url"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"

)


func dataSourceResgroupRead(d *schema.ResourceData, m interface{}) error {
	controller := m.(*ControllerCfg)

	name := d.Get("name").(string)
	tenant := d.Get("tenant").(string)

	url_values := &url.Values{}
	resp, err := controller.decsAPICall("POST", "/restmachine/cloudapi/cloudspaces/list", url_values)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	model := CloudspacesListResp{}
	json.NewDecoder(resp.Body).Decode(&model)

	// d.SetId(model[matching_index].id)
	// d.Set("field_name", value)
	return nil
}


func dataSourceResgroup() *schema.Resource {
	return &schema.Resource {
		SchemaVersion: 1,

		Read:   dataSourceResgroupRead,

		Timeouts: &schema.ResourceTimeout {
			Read:    &Timeout30s,
			Default: &Timeout60s,
		},

		Schema: map[string]*schema.Schema {
			"name": {
				Type:          schema.TypeString,
				Required:      true,
				Description:  "Name of this resource group. Names are unique within the context of a tenant and case sensitive.",
			},

			"tenant": &schema.Schema {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the tenant, which this resource group belongs to.",
			},
		},
	}
}