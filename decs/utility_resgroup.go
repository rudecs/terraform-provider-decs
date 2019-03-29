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
	// "strconv"

	// "github.com/hashicorp/terraform/helper/schema"
	// "github.com/hashicorp/terraform/helper/validation"
)

func (ctrl *ControllerCfg) utilityResgroupConfigGet(rgid int) (*ResgroupConfig, error) {
	url_values := &url.Values{}
	url_values.Add("cloudspaceId", fmt.Sprintf("%d", rgid))
	resgroup_facts, err := ctrl.decsAPICall("POST", CloudspacesGetAPI, url_values)
	if err != nil {
		return nil, err
	}

	log.Printf("utilityResgroupConfigGet: ready to unmarshal string %q", resgroup_facts)
	model := CloudspacesGetResp{}
	err = json.Unmarshal([]byte(resgroup_facts), &model)
	if err != nil {
		return nil, err
	}

	ret := &ResgroupConfig{}
	ret.TenantID = model.TenantID
	ret.Location = model.Location
	ret.Name = model.Name
	ret.ID = rgid
	ret.GridID = model.GridID
	ret.ExtIP = model.ExtIP   // legacy field for VDC - this will eventually become obsoleted by true Resource Groups
	// Quota ResgroupQuotaConfig
	// Network NetworkConfig
	log.Printf("utilityResgroupConfigGet: tenant ID %d, GridID %d, ExtIP %q", 
	           model.TenantID, model.GridID, model.ExtIP)

	return ret, nil
}