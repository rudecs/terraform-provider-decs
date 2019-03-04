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

)

func Jo2JSON(arg_str string) string {
	ret_string := strings.Replace(string(arg_str), "u'", "\"", -1)
	ret_string = strings.Replace(ret_string, "'", "\"", -1)
	ret_string = strings.Replace(ret_string, "False", "false", -1)
	ret_string = strings.Replace(ret_string, "True", "true", -1)
	return ret_string
}
