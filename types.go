package main

type Workload struct {
	WorkLoadName string `json:"workLoadName"`
	BackupType string `json:"backupType"`
	Site string `json:"site"`
	//size in TB
	WorkLoadCap int64 `json:"workLoadCap"`
	ProcessCap int64 `json:"processCap"`
	WorkLoadCapBytes int64 `json:"workLoadCapBytes"`
	ProcessCapBytes int64 `json:"processCapBytes"`
	VMQuantity int `json:"vmQty"`
	VMVMDKRatio int `json:"vmVmdkRatio"`
	VMVMDKQuantity int `json:"vmdkQty"`
}


func NewWorkLoad() (Workload) {
	w := Workload{
		"Default",
		"vm",
		"Site_A",
        0,
		0,
		0,
		0,
		0,
		0,
		0,
	}
	return w
}

type VSEFormat struct {
	Workloads []Workload `json:"workload"`
}


type ShadowPart struct {
	PartName string
	Free int64
	Capacity int64
	Used int64
}
type ShadowDisk struct {
	DiskName string
	Alloc int64
}
type ShadowAssigned struct {
	WorkLoadCapGB int64
	ProcessCapGB int64
}
type ShadowVM struct {
	Name string
	GuestOS string
	IsWin bool
	Disks []ShadowDisk
	Parts []ShadowPart
	DiskAllocGB int64
	PartSizeGB int64
	Assigned ShadowAssigned
	ComputeResource string
}
type ShadowGroup struct {
	VMs []ShadowVM
	GroupName string
}
type ShadowFile struct {
	GroupedVM []ShadowGroup
}