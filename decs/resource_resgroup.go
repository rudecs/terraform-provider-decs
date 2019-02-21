/*
Copyright (c) 2019 Digital Energy Cloud Solutions. All Rights Reserved.

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
)

func resourceResgroupCreate(d *schema.ResourceData, m interface{}) error {
	return resourceResgroupRead(d, m)
}

func resourceResgroupRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceResgroupUpdate(d *schema.ResourceData, m interface{}) error {
	return resourceResgroupRead(d, m)
}

func resourceResgroupDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceResgroup() *schema.Resource {
	return &schema.Resource {
		Create: resourceResgroupCreate,
		Read:   resourceResgroupRead,
		Update: resourceResgroupUpdate,
		Delete: resourceResgroupDelete,

		Schema: map[string]*schema.Schema {
			"name": &schema.Schema {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}