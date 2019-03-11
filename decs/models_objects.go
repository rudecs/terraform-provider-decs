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

type DiskConfig struct {
	Label string
	Size int
	Pool string
	Provider string
}

type NetworkConfig struct {
	Label string
	NetworkID int
}

type PortforwardConfig struct {
	Label string
	ExtPort int
	IntPort int
	Proto string
}

type SshKeyConfig struct {
	User string
	SshKey string
}

type MachineConfig struct {
	ResGroupID int
	Name string
	ID int
	Cpu int
	Ram int
	ImageID int
	BootDisk DiskConfig
	DataDisks []DiskConfig
	Networks []NetworkConfig
	PortForwards []PortforwardConfig
	SshKeys []SshKeyConfig
	Description string
}

type ResgroupQuotaConfig struct {
	Cpu int
	Ram int
	Disk int
	NetTraffic int
	ExtIPs int
}

type ResgroupConfig struct {
	TenantID int
	Location string
	Name string
	ID int
	Quota ResgroupQuotaConfig
	Network NetworkConfig
}