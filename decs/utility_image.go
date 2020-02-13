/*
Copyright (c) 2020 Digital Energy Cloud Solutions LLC. All Rights Reserved.
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
	"net/url"
	"strconv"

	"github.com/hashicorp/terraform/helper/schema"
	// "github.com/hashicorp/terraform/helper/validation"
)

func (ctrl *ControllerCfg) utilityGetImageByID(imgrec *ImageRecord) error {
	// This method will try to find OS image by the ID passed in imgrec.ID and read
	// in all other image details. The details will be written back into imgrec
	// structure.
	//
	// Intended use case for this function is when creating disks for a VM without
	// explicitly telling which SEP and pool to use. Then by default SEP ID and pool 
	// name will be inherited from those of the OS image used to spin off the VM.
	// As OS image is passed to VM provider via its ID, we need this convenience
	// method to get SEP ID & pool.

	return nil
}