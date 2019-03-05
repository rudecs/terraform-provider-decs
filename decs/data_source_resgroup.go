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

	"github.com/hashicorp/terraform/helper/schema"
	// "github.com/hashicorp/terraform/helper/validation"

)


func dataSourceResgroupRead(d *schema.ResourceData, m interface{}) error {
	name := d.Get("name").(string)
	tenant_name := d.Get("tenant").(string)

	controller := m.(*ControllerCfg)
	url_values := &url.Values{}
	body_string, err := controller.decsAPICall("POST", CloudspacesListAPI, url_values)
	if err != nil {
		return err
	}

	log.Printf("%s", body_string)
	log.Printf("dataSourceResgroupRead: ready to decode response body")
	model := CloudspacesListResp{}
	err = json.Unmarshal([]byte(body_string), &model)
	if err != nil {
		return err
	}

	/*
	log.Printf("dataSourceResgroupRead: ready to decode response body")
	err = json.NewDecoder(resp.Body).Decode(&model)
	if err != nil {
		return err
	}

	
	// var model []CloudspaceRecord
	model := CloudspacesListResp{}
	name = "vdc01"
	tenant_name = "GreyseTmpFromJS"
	JsonPart := `[{"status":"DEPLOYED","updateTime":1523027184,"externalnetworkip":"185.193.143.152","name":"vdc01","descr":"","creationTime":1523027135,"acl":[{"status":"CONFIRMED","canBeDeleted":false,"right":"ACDRUX","type":"U","userGroupId":"vadim_sorokin_1@itsyouonline"},{"status":"CONFIRMED","canBeDeleted":false,"right":"CRX","type":"U","userGroupId":"svs1370g@itsyouonline"}],"accountAcl":{"status":"CONFIRMED","right":"RCX","explicit":true,"userGroupId":"svs1370g@itsyouonline","guid":"","type":"U"},"gid":2001,"location":"ds1","publicipaddress":"185.193.143.152","accountName":"GreyseTmpFromJS","id":76,"accountId":21},{"status":"DEPLOYED","updateTime":1523027184,"externalnetworkip":"185.193.143.152","name":"vdc01","descr":"","creationTime":1523027135,"acl":[{"status":"CONFIRMED","canBeDeleted":false,"right":"ACDRUX","type":"U","userGroupId":"vadim_sorokin_1@itsyouonline"},{"status":"CONFIRMED","canBeDeleted":false,"right":"CRX","type":"U","userGroupId":"svs1370g@itsyouonline"}],"accountAcl":{"status":"CONFIRMED","right":"RCX","explicit":true,"userGroupId":"svs1370g@itsyouonline","guid":"","type":"U"},"gid":2001,"location":"ds1","publicipaddress":"185.193.143.152","accountName":"GreyseTmpFromJS","id":76,"accountId":21}]`
	err := json.Unmarshal([]byte(JsonPart), &model)
	if err != nil {
		return err
	}
	*/

	log.Printf("dataSourceResgroupRead: traversing decoded Json of length %d", len(model))
	for index, item := range model {
		// need to match VDC by name & tenant name
		if item.Name == name && item.TenantName == tenant_name {
			log.Printf("dataSourceResgroupRead: index %d, name %q, tenant %q", index, item.Name, item.TenantName)
			d.SetId(fmt.Sprintf("%d", model[index].ID))
			d.Set("tenant_id", model[index].TenantID)
			// d.Set("field_name", value)
			return nil
		}
	}

	return fmt.Errorf("Cannot find resource group name %q owned by tenant %q", name, tenant_name)
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

			"tenant_id": &schema.Schema {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Unique ID of the tenant, which this resource group belongs to.",
			},

			"public_ip": {
				Type:          schema.TypeString,
				Computed:      true,
				Description:  "Public IP address of this resource group (if any).",
			},
		},
	}
}