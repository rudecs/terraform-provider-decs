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

	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	// "github.com/hashicorp/terraform/terraform"

)

var decsController *ControllerCfg

func Provider() *schema.Provider {
	return &schema.Provider {
		Schema: map[string]*schema.Schema {
			"authenticator": {
				Type:        schema.TypeString,
				Required:    true,
				StateFunc:   stateFuncToLower,
				ValidateFunc: validation.StringInSlice([]string{"oauth2", "legacy", "jwt"}, true), // ignore case while validating
				Description: "Authentication mode to use when connecting to DECS cloud API. Should be one of 'oauth2', 'legacy' or 'jwt'.",
			},

			"oauth2_url": {
				Type:        schema.TypeString,
				Optional:    true,
				StateFunc:   stateFuncToLower,
				DefaultFunc: schema.EnvDefaultFunc("DECS_OAUTH2_URL", nil),
				Description: "The Oauth2 application URL in 'oauth2' authentication mode.",
			},

			"controller_url": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				StateFunc:   stateFuncToLower,
				Description: "The URL of DECS Cloud controller to use. API calls will be directed to this URL.",
			},

			"user": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("DECS_USER", nil),
				Description: "The user name for DECS cloud API operations in 'legacy' authentication mode.",
			},

			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("DECS_PASSWORD", nil),
				Description: "The user password for DECS cloud API operations in 'legacy' authentication mode.",
			},

			"app_id": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("DECS_APP_ID", nil),
				Description: "Application ID to access DECS cloud API in 'oauth2' authentication mode.",
			},

			"app_secret": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("DECS_APP_SECRET", nil),
				Description: "Application secret to access DECS cloud API in 'oauth2' authentication mode.",
			},

			"jwt": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("DECS_JWT", nil),
				Description: "JWT to access DECS cloud API in 'jwt' authentication mode.",
			},

			"allow_unverified_ssl": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "If set, DECS API will allow unverifiable SSL certificates.",
			},
		},
		
		ResourcesMap: map[string]*schema.Resource {
			"decs_resgroup": resourceResgroup(),
			"decs_vm": resourceVm(),
		},

		DataSourcesMap: map[string]*schema.Resource {
			"decs_resgroup": dataSourceResgroup(),
			"decs_vm": dataSourceVm(),
			"decs_image": dataSourceImage(),
		},
		
		ConfigureFunc: providerConfigure,
	}
}

func stateFuncToLower(argval interface{}) string {
	return strings.ToLower(argval.(string))
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	decsController, err := ControllerConfigure(d)
	if err != nil {
		return nil, err
	}
	return decsController, nil
}