package main

type Workload struct {
	WorkLoadName string `json:"workLoadName"`
	BackupType   string `json:"backupType"`
	Site         string `json:"site"`
	CopySite     string `json:"copySite"`
	//size in TB
	WorkLoadCap      int64 `json:"workLoadCap"`
	ProcessCap       int64 `json:"processCap"`
	WorkLoadCapBytes int64 `json:"workLoadCapBytes"`
	ProcessCapBytes  int64 `json:"processCapBytes"`
	VMQuantity       int   `json:"vmQty"`
	VMVMDKRatio      int   `json:"vmVmdkRatio"`
	VMVMDKQuantity   int   `json:"vmdkQty"`

	GrowthPercent int     `json:"growthPercent"`
	ScopeYears    int     `json:"scopeYears"`
	BackupWindow  int     `json:"backupWindow"`
	ChangeRate    int     `json:"changeRate"`
	Reduction     int     `json:"reduction"`
	UsePerVM      string  `json:"usePerVM"`
	UseReFs       string  `json:"useReFs"`
	RpsBu         int     `json:"rpsBu"`
	BuWeekly      int     `json:"buWeekly"`
	BuMonthly     int     `json:"buMonthly"`
	BuYearly      int     `json:"buYearly"`
	RpsBuCopy     int     `json:"rpsBuCopy"`
	BuCopyWeekly  int     `json:"buCopyWeekly"`
	BuCopyMonthly int     `json:"buCopyMonthly"`
	BuCopyYearly  int     `json:"buCopyYearly"`
	CloudMove     int     `json:"cloudMove"`
	CloudEnabled  bool    `json:"cloudEnabled"`
	BandWidthInc  float32 `json:"bandWidthInc"`
}

func NewWorkLoad() Workload {
	w := Workload{
		"Default",
		"vm",
		"Site_A",
		"Site_B",
		0,
		0,
		0,
		0,
		0,
		0,
		0,
		//defaults
		10,
		3,
		8,
		5,
		50,
		"perVM",
		"yes",
		30,
		1,
		0,
		0,
		30,
		1,
		0,
		0,
		0,
		false,
		0,
	}
	return w
}

type VSEFormat struct {
	Workloads []Workload `json:"workload"`
}

type ShadowPart struct {
	PartName string
	Free     int64
	Capacity int64
	Used     int64
}
type ShadowDisk struct {
	DiskName string
	Alloc    int64
}
type ShadowAssigned struct {
	WorkLoadCapGB int64
	ProcessCapGB  int64
}
type ShadowVM struct {
	Name            string
	GuestOS         string
	IsWin           bool
	Disks           []ShadowDisk
	Parts           []ShadowPart
	DiskAllocGB     int64
	PartSizeGB      int64
	Assigned        ShadowAssigned
	ComputeResource string
}
type ShadowGroup struct {
	VMs       []ShadowVM
	GroupName string
}
type ShadowFile struct {
	GroupedVM []ShadowGroup
}
