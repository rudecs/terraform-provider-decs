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

	"time"
)

//
// timeouts for API calls from CRUD functions of Terraform plugin
var Timeout30s = time.Second * 30
var Timeout60s = time.Second * 60
var Timeout180s = time.Second * 180

//
// structures related to /cloudapi/cloudspaces/list API
//
type UserAclRecord struct {
	status string          `json:"status"`
	can_delete bool        `json:"canBeDeleted"`
	acc_rights string      `json:"right"`
	acc_type string        `json:"type"`
	ugroup_id string       `json:"userGroupId"`
}

type AccountAclRecord struct {
	status string          `json:"status"`
	explicit bool          `json:"explicit"`
	acc_rights string      `json:"right"`
	acc_type string        `json:"type"`
	entity_id string       `json:"userGroupId"`
	guid string            `json:"guid"`
}

type CloudspaceRecord struct {
	status string          `json:"status"`
	update_time uint64     `json:"updateTime"`
	ext_net_ip string      `json:"externalnetworkip"`
	name string            `json:"name"`
	decsription string     `json:"description"`
	create_time uint64     `json:"creationTime"`
	acl []UserAclRecord    `json:"acl"`
	owner AccountAclRecord `json:"accountAcl"`
	grid_id int            `json:"gid"`
	location string        `json:"location"`
	public_ip string       `json:"publicipaddress"`
	account_name string    `json:"accountName"`
	id uint                `json:"id"`
	account_id int         `json:"accountId"`
}

const CloudspacesListAPI = "/restmachine/cloudapi/cloudspaces/list"
type CloudspacesListResp []CloudspaceRecord

//
// structures related to /cloudapi/cloudspaces/create API call
//
const CloudspacesCreateAPI= "/restmachine/cloudapi/cloudspaces/create"
type CloudspacesCreateParam struct {
	account_id int         `json:"accountId"`
	location string        `json:"location"`
	name string            `json:"name"`
	owner string           `json:"access"`
	cpu int                `json:"maxCPUCapacity"`
	ram int                `json:"maxMemoryCapacity"`
	disk int               `json:"maxVDiskCapacity"`
	net_traffic int        `json:"maxNetworkPeerTransfer"`
	ext_ips int            `json:"maxNumPublicIP"`
	ext_net_id int         `json:"externalnetworkid"`
	allowed_size_ids []int `json:"allowedVMSizes"`
	int_net_range string   `json:"privatenetwork"`
} 

//
// structures related to /cloudapi/cloudspaces/get API call
//
type QuotaRecord struct {
	cpu int                `json:"CU_C"`
	ram int                `json:"CU_M"`
	disk int               `json:"CU_D"`
	net_traffic int        `json:"CU_NP"`
	ext_ips int            `json:"CU_I"`
}

const CloudspacesGetAPI= "/restmachine/cloudapi/cloudspaces/get"
type CloudspacesGetResp struct {
	status string          `json:"status"`
	update_time uint64     `json:"updateTime"`
	ext_net_ip string      `json:"externalnetworkip"`
	decsription string     `json:"description"`
	quotas QuotaRecord     `json:"resourceLimits"`
	id uint                `json:"id"`
	account_id int         `json:"accountId"`
	name string            `json:"name"`
	create_time uint64     `json:"creationTime"`
	acl []UserAclRecord    `json:"acl"`
	secret string          `json:"secret"`
	grid_id int            `json:"gid"`
	location string        `json:"location"`
	public_ip string       `json:"publicipaddress"`
}

// 
// structures related to /cloudapi/cloudspaces/update API
//
const CloudspacesUpdateAPI = "/restmachine/cloudapi/cloudspaces/update"
type CloudspacesUpdateParam struct {
	id uint                `json:"cloudspaceId"`
	name string            `json:"name"`
	cpu int                `json:"maxCPUCapacity"`
	ram int                `json:"maxMemoryCapacity"`
	disk int               `json:"maxVDiskCapacity"`
	net_traffic int        `json:"maxNetworkPeerTransfer"`
	ext_ips int            `json:"maxNumPublicIP"`
}

//
// structures related to /cloudapi/machines/create API
//
const MachineCreateAPI = "/restmachine/cloudapi/machines/create"
type MachineCreateParam struct {
	vdc_id uint            `json:"cloudspaceId"`
	name string            `json:"name"`
	description string     `json:"description"`
	cpu int                `json:"vcpus"`
	ram int                `json:"memory"`
	image_id int           `json:"imageId"`
	boot_disk int          `json:"disksize"`
	data_disks []int       `json:"datadisks"`
	user_data string       `json:"userdata"`
}

// 
// structures related to /cloudapi/machines/list API
//
type NicRecord struct {
	status string          `json:"status"`
	mac_address string     `json:"macAddress"`
	reference_id string    `json:"referenceId"`
	device_name string     `json:"deviceName"`
	nic_type string        `json:"type"`
	params string          `json:"params"`
	network_id int         `json:"networkId"`
	guid string            `json:"guid"`
	ip_address string      `json:"ipAddress"`
}

type MachineRecord struct {
	status string          `json:"status"`
	stack_id int           `json:"stackId"`
	update_time uint64     `json:"updateTime"`
	reference_id string    `json:"referenceId"`
	name string            `json:"name"`
	nics []NicRecord       `json:"nics"`
	size_id int            `json:"sizeId"`
	data_disks []uint      `json:"disks"`
	create_time uint64     `json:"creationTime"`
	image_id int           `json:"imageId"`
	boot_disk int          `json:"storage"`
	cpu int                `json:"vcpus"`
	ram int                `json:"memory"`
	id uint                `json:"id"`
}

const MachinesListAPI = "/restmachine/cloudapi/machines/list"
type MachinesListResp []MachineRecord

//
// structures related to /cloudapi/machines/get
//
type DataDiskRecord struct {
	status string          `json:"status"`
	size_max int           `json:"sizeMax"`
	label string           `json:"name"`
	description string     `json:"descr"`
	acl map[string]string  `json:"acl"`
	disk_type string       `json:"type"`
	id uint                `json:"id"`
}

type GuestLoginRecord struct {
	guid string            `json:"guid"`
	login string           `json:"login"`
	password string        `json:"password"`
}

const MachinesGetAPI = "/restmachine/cloudapi/machines/get"
type MachinesGetResp struct {
	vdc_id uint            `json:"cloudspaceId`
	status string          `json:"status"`
	update_time uint64     `json:"updateTime"`
	hostname string        `json:"hostname"`
	is_locked bool         `json:"locked"`
	name string            `json:"name"`
	create_time uint64     `json:"creationTime"`
	size_id uint           `json:"sizeid"`
	cpu int                `json:"vcpus"`
	ram int                `json:"memory"`
	boot_disk int          `json:"storage"`
	data_disks []DiskRecord `json:"disks"`
	nics []NicRecord       `json:"interfaces"`
	guest_logins []GuestLoginRecord `json:"accounts"`
	image_name string      `json:"osImage"`
	image_id int           `json:"imageid"`
	description string     `json:"description"`
	id uint                `json:"id"`
}

//
// structures related to /restmachine/cloudapi/images/list API
//
type ImagesRecord struct {
	status string       `json:"status"`
	username string     `json:"username"`
	description string  `json:"description"`
	account_id uint     `json:"accountId"`
	size int            `json:"size"`
	image_type string   `json:"type"`
	id uint             `json:"id"`
	name string         `json:"name"`
}

const ImagesListAPI = "/restmachine/cloudapi/images/list"
type ImagesListResp []ImageRecord

//
// structures related to /cloudapi/externalnetwork/list API
//
type ExtNetworkRecord struct {
	ip_range string        `json:"name"`
	id uint                `json:"id"`
} 

const  ExtNetworksListAPI = "/restmachine/cloudapi/externalnetwork/list"
type ExtNetworksResp []ExtNetRecord

//
// structures related to /cloudapi/accounts/list API
//
type TenantRecord struct {
	id int                 `json:"id"`
	update_time uint64     `json:"updateTime"`
	create_time uint64     `json:"creationTime"`
	name string            `json:"name"`
	acl []UserAclRecord    `json:"acl"`
}

const TenantsListAPI = "/restmachine/cloudapi/accounts/list"
type TenantsListResp []TenantRecord