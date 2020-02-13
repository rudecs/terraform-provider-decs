/*
Copyright (c) 2019-2020 Digital Energy Cloud Solutions LLC. All Rights Reserved.
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
	Status string          `json:"status"`
	CanBeDeleted bool      `json:"canBeDeleted"`
	AccRights string       `json:"right"`
	AccType string         `json:"type"`
	UgroupID string        `json:"userGroupId"`
}

type AccountAclRecord struct {
	Status string          `json:"status"`
	AccRights string       `json:"right"`
	IsExplicit bool        `json:"explicit"`
	EntityID string        `json:"userGroupId"`
	Guid string            `json:"guid"`
	AccType string         `json:"type"`
}

type CloudspaceRecord struct {
	Status string          `json:"status"`
	UpdateTime uint64      `json:"updateTime"`
	ExtNetIP string        `json:"externalnetworkip"`
	Name string            `json:"name"`
	Decsription string     `json:"descr"`
	CreateTime uint64      `json:"creationTime"`
	Acl []UserAclRecord    `json:"acl"`
	Owner AccountAclRecord `json:"accountAcl"`
	GridID int             `json:"gid"`
	Location string        `json:"location"`
	PublicIP string        `json:"publicipaddress"`
	TenantName string      `json:"accountName"`
	ID uint                `json:"id"`
	TenantID int           `json:"accountId"`
}

const CloudspacesListAPI = "/restmachine/cloudapi/cloudspaces/list"
type CloudspacesListResp []CloudspaceRecord

//
// structures related to /cloudapi/cloudspaces/create API call
//
const ResgroupCreateAPI= "/restmachine/cloudapi/cloudspaces/create"
type ResgroupCreateParam struct {
	TenantID int           `json:"accountId"`
	Location string        `json:"location"`
	Name string            `json:"name"`
	Owner string           `json:"access"`
	Cpu int                `json:"maxCPUCapacity"`
	Ram int                `json:"maxMemoryCapacity"`
	Disk int               `json:"maxVDiskCapacity"`
	NetTraffic int         `json:"maxNetworkPeerTransfer"`
	ExtIPs int             `json:"maxNumPublicIP"`
	ExtNetID int           `json:"externalnetworkid"`
	AllowedSizeIDs []int   `json:"allowedVMSizes"`
	IntNetRange string     `json:"privatenetwork"`
} 

//
// structures related to /cloudapi/cloudspaces/update API call
//
const ResgroupUpdateAPI= "/restmachine/cloudapi/cloudspaces/update"

//
// structures related to /cloudapi/cloudspaces/get API call
//
type QuotaRecord struct {
	Cpu int                `json:"CU_C"`
	Ram float32            `json:"CU_M"` // NOTE: it is float32! Casting to int may be required when passing it to ResgroupConfig
	Disk int               `json:"CU_D"`
	NetTraffic int         `json:"CU_NP"`
	ExtIPs int             `json:"CU_I"`
}

const CloudspacesGetAPI= "/restmachine/cloudapi/cloudspaces/get"
type CloudspacesGetResp struct {
	Status string          `json:"status"`
	UpdateTime uint64      `json:"updateTime"`
	ExtIP string           `json:"externalnetworkip"`
	Description string     `json:"description"`
	Quotas QuotaRecord     `json:"resourceLimits"`
	ID uint                `json:"id"`
	TenantID int           `json:"accountId"`
	Name string            `json:"name"`
	CreateTime uint64      `json:"creationTime"`
	Acl []UserAclRecord    `json:"acl"`
	Secret string          `json:"secret"`
	GridID int             `json:"gid"`
	Location string        `json:"location"`
	PublicIP string        `json:"publicipaddress"`
}

// 
// structures related to /cloudapi/cloudspaces/update API
//
const CloudspacesUpdateAPI = "/restmachine/cloudapi/cloudspaces/update"
type CloudspacesUpdateParam struct {
	ID uint                `json:"cloudspaceId"`
	Name string            `json:"name"`
	Cpu int                `json:"maxCPUCapacity"`
	Ram int                `json:"maxMemoryCapacity"`
	Disk int               `json:"maxVDiskCapacity"`
	NetTraffic int         `json:"maxNetworkPeerTransfer"`
	ExtIPs int             `json:"maxNumPublicIP"`
}

// 
// structures related to /cloudapi/cloudspaces/delete API
//
const CloudspacesDeleteAPI = "/restmachine/cloudapi/cloudspaces/delete"

//
// structures related to /cloudapi/machines/create API
//
const MachineCreateAPI = "/restmachine/cloudapi/machines/create"
type MachineCreateParam struct {
	ResGroupID uint        `json:"cloudspaceId"`
	Name string            `json:"name"`
	Description string     `json:"description"`
	Cpu int                `json:"vcpus"`
	Ram int                `json:"memory"`
	ImageID int            `json:"imageId"`
	BootDisk int           `json:"disksize"`
	DataDisks []int        `json:"datadisks"`
	UserData string        `json:"userdata"`
}

// strucures related to cloudapi/machines/delete API
const MachineDeleteAPI = "/restmachine/cloudapi/machines/delete"

// 
// structures related to /cloudapi/machines/list API
//
type NicRecord struct {
	Status string          `json:"status"`       // did not see any other values but ""
	MacAddress string      `json:"macAddress"`   // example "52:54:00:00:2d:2a"
	ReferenceID string     `json:"referenceId"`  // did not see any other values but ""
	DeviceName string      `json:"deviceName"`   // internal: "vm-13578-00a0", external: "vm-13578-6-ext
	NicType string         `json:"type"`         // "bridge" for int net, "PUBLIC" for ext net
	Params string          `json:"params"`       // for ext net "gateway:176.118.165.1 externalnetworkId:6"
	NetworkID int          `json:"networkId"`    // did not see any other values but 0
	Guid string            `json:"guid"`         // did not see any other values but ""
	IPAddress string       `json:"ipAddress"`    // example "176.118.165.25/24"
}

type MachineRecord struct {
	Status string          `json:"status"`
	StackID int            `json:"stackId"`
	UpdateTime uint64      `json:"updateTime"`
	ReferenceID string     `json:"referenceId"`
	Name string            `json:"name"`
	NICs []NicRecord       `json:"nics"`
	SizeID int             `json:"sizeId"`
	DataDisks []uint       `json:"disks"`
	CreateTime uint64      `json:"creationTime"`
	ImageID int            `json:"imageId"`
	BootDisk int           `json:"storage"`
	Cpu int                `json:"vcpus"`
	Ram int                `json:"memory"`
	ID uint                `json:"id"`
}

const MachinesListAPI = "/restmachine/cloudapi/machines/list"
type MachinesListResp []MachineRecord

//
// structures related to /cloudapi/machines/get
//
type GenericDiskRecord struct {
	Status string          `json:"status"`
	SizeMax int            `json:"sizeMax"`
	Label string           `json:"name"`
	Description string     `json:"descr"`
	Acl map[string]string  `json:"acl"`
	Type string            `json:"type"`
	ID uint                `json:"id"`
	SepId int              `json:"sepid"`
	Pool string            `json:"pool"`
	AccountId int          `json:"accountId"`
	VmId int               `json:"vmid"`
	ParentId uint          `json:"parentId"`
	TechStatus string      `json:"techStatus"`
	ImageId int            `json:"imageId"`
}

type GuestLoginRecord struct {
	Guid string            `json:"guid"`
	Login string           `json:"login"`
	Password string        `json:"password"`
}

const MachinesGetAPI = "/restmachine/cloudapi/machines/get"
type MachinesGetResp struct {
	ResGroupID uint           `json:"cloudspaceid"` // note that "id" is not capitalized in "cloudspaceid"
	Status string             `json:"status"`
	UpdateTime uint64         `json:"updateTime"`
	Hostname string           `json:"hostname"`
	IsLocked bool             `json:"locked"`
	Name string               `json:"name"`
	CreateTime uint64         `json:"creationTime"`
	SizeID uint               `json:"sizeid"`
	Cpu int                   `json:"vcpus"`
	Ram int                   `json:"memory"`
	BootDisk int              `json:"storage"` // this is requested boot disk size in GB, not boot disk ID
	Disks []GenericDiskRecord `json:"disks"`   // all disks associated with this VM, boot & data
	NICs []NicRecord          `json:"interfaces"`
	GuestLogins []GuestLoginRecord `json:"accounts"`
	ImageName string          `json:"osImage"`
	ImageID int               `json:"imageid"`
	Description string        `json:"description"`
	ID uint                   `json:"id"`
}

//
// structures related to /restmachine/cloudapi/images/list API
//
type ImageRecord struct {
	Status string       `json:"status"`
	Username string     `json:"username"`
	Description string  `json:"description"`
	TenantID uint       `json:"accountId"`
	Size int            `json:"size"`
	Type string         `json:"type"`
	ID uint             `json:"id"`
	Name string         `json:"name"`
	SepId int           `json:"sepid"`
	Pool string         `json:"pool"`
}

const ImagesListAPI = "/restmachine/cloudapi/images/list"
type ImagesListResp []ImageRecord

//
// structures related to /cloudapi/externalnetwork/list API
//
type ExtNetworkRecord struct {
	IPRange string       `json:"name"`
	ID uint              `json:"id"`
} 

const AccountExtNetworksListAPI = "/restmachine/cloudapi/externalnetwork/list"
type AcountExtNetworksResp []ExtNetworkRecord

//
// Response of this call in current API version is just a list of attached network IDs
const VmExtNetworksListAPI = "/restmachine/cloudapi/machines/listExternalNetworks"

//
// structures related to /cloudapi/accounts/list API
//
type TenantRecord struct {
	ID int                 `json:"id"`
	UpdateTime uint64      `json:"updateTime"`
	CreateTime uint64      `json:"creationTime"`
	Name string            `json:"name"`
	Acl []UserAclRecord    `json:"acl"`
}

const TenantsListAPI = "/restmachine/cloudapi/accounts/list"
type TenantsListResp []TenantRecord

//
// structures related to /cloudapi/portforwarding/list API
//
type PortforwardRecord struct {
	Proto string           `json:"protocol"`
	IntPort string         `json:"localPort"`
	ExtPort string         `json:"publicPort"`
	ExtIP string           `json:"publicIp"`
	IntIP string           `json:"localIp"`
	VmID int               `json:"machineId"`
	VmName string          `json:"machineName"`
} 

const PortforwardsListAPI = "/restmachine/cloudapi/portforwarding/list"
type PortforwardsResp []PortforwardRecord

const PortforwardingCreateAPI = "/restmachine/cloudapi/portforwarding/create"

//
// structures related to /cloudapi/machines/attachExternalNetwork API
//
const AttachExternalNetworkAPI = "/restmachine/cloudapi/machines/attachExternalNetwork"

//
//
//
const DiskCreateAPI = "/restmachine/cloudapi/disks/create"
const DiskAttachAPI = "/restmachine/cloudapi/machines/attachDisk"
