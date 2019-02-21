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

	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"

)

// enumerated constants that define authentication modes 
const (
	MODE_UNDEF   = iota // this is the invalid mode - it should never be seen
	MODE_LEGACY  = iota
	MODE_OAUTH2  = iota
	MODE_JWT     = iota
)

type ControllerCfg struct {
	controller_url   string  // always required
	auth_mode        int     // always required
	auth_mode_txt    string  // always required, it is a text representation of auth mode
	legacy_user      string  // required for legacy mode
	legacy_password  string  // required for legacy mode
	legacy_sid       string  // obtained from DECS controller on successful login in legacy mode
	jwt              string  // obtained from Outh2 provider on successful login in oauth2 mode, required in jwt mode
	app_id           string  // required for oauth2 mode
	app_secret       string  // required for oauth2 mode
	oauth2_url       string  // always required
	oauth2_username  string  // obtained from Outh2 provider on successful login in oauth2 or jwt modes
}

func ControllerConfig(d *schema.ResourceData) (*ControllerCfg, error) {
	// This function will check that all required provider parameters for the 
	// selected authenticator mode are set correctly and initialize ControllerCfg structure
	// based on the provided parameters.
	// It will NOT check for validity of supplied credentials - this is done by ControllerClient method
	// that actually initiates connection to the specified DECS controller URL and, if succeeded, completes 
	// ControllerCfg structure with the rest of computed parameters (e.g. JWT, session ID and Oauth2 user name).
	//
	// The structure created by this function should be used with ControllerClient method to initalize connection
	// to the DECS cloud controller.

	ret_config := &ControllerCfg{
		controller_url:  d.Get("controller_url").(string),
		auth_mode:       MODE_UNDEF,
		legacy_user:     d.Get("user").(string),
		legacy_password: d.Get("password").(string),
		legacy_sid:      "",
		jwt:             d.Get("jwt").(string),
		app_id:          d.Get("app_id").(string),
		app_secret:      d.Get("app_secret").(string),
		oauth2_url:      d.Get("oauth2_url").(string),
		oauth2_username: "",
	}

	if ret_config.controller_url = "" {
		return nil, fmt.Errorf("Empty DECS cloud controller URL provided.")
	}

	ret_config.auth_mode_txt = strings.ToLower(d.Get("authenticator").(string))
	switch ret_config.auth_mode_txt {
	case "jwt":
		if ret_config.jwt = "" {
			return nil, fmt.Errorf("Authenticator mode 'jwt' specified but no JWT provided.")
		}
		ret_config.auth_mode = MODE_JWT
	case "oauth2":
		if ret_config.oauth2_url = "" {
			return nil, fmt.Errorf("Authenticator mode 'oauth2' specified but no OAuth2 URL provided.")
		}
		if ret_config.app_id = "" {
			return nil, fmt.Errorf("Authenticator mode 'oauth2' specified but no Application ID provided.")
		}
		if ret_config.app_secret = "" {
			return nil, fmt.Errorf("Authenticator mode 'oauth2' specified but no Secret ID provided.")
		}
		ret_config.auth_mode = MODE_OAUTH2
	case "legacy":
		//
		ret_config.user = d.Get("user").(string)
		if ret_config.user = "" {
			return nil, fmt.Errorf("Authenticator mode 'legacy' specified but no user provided.")
		}
		ret_config.password = d.Get("password").(string)
		if ret_config.password = "" {
			return nil, fmt.Errorf("Authenticator mode 'legacy' specified but no password provided.")
		}
		ret_config.auth_mode = MODE_LEGACY
	default:
		return nil, fmt.Errorf("Unknown authenticator mode '" + ret_config.auth_mode_txt + "' provided.")
	}

	return ret_cfg, nil
}

func (config *ControllerCfg) ControllerClient() error {
}

func (config *ControllerCfg) getOAuth2JWT() (string, error) {
	if config.auth_mode == MODE_UNDEF {
		return nil, fmt.Errorf("getOAuth2JWT method called for undefined authorization mode.")
	}
	if config.auth_mode != MODE_OAUTH2 {
		return nil, fmt.Errorf("getOAuth2JWT method called for incompatible authorization mode '" + config.auth_mode_txt + "'.")
	}

	return config.jwt, nil
}

func (config *ControllerCfg) validateJWT(jwt *string) (bool, error) {
	if jwt == "" {
		if config.jwt == "" {
			return false, fmt.Errorf("validateJWT method called, but no meaningful JWT is available.")
		}
		jwt = config.jwt
	}

	return true, nil
}

func (config *ControllerCfg) validateLegacyUser() (bool, error) {
	if config.auth_mode == MODE_UNDEF {
		return nil, fmt.Errorf("validateLegacyUser method called for undefined authorization mode.")
	}
	if config.auth_mode != MODE_LEGACY {
		return nil, fmt.Errorf("validateLegacyUser method called for incompatible authorization mode '" + config.auth_mode_txt + "'.")
	}
	
	return true, nil
}

func (config *ControllerCfg) decsAPICall() error {
	
}

