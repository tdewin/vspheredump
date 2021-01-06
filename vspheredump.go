package main

import (
	"context"
	"encoding/json"
	"fmt"

	"io/ioutil"

	"github.com/vmware/govmomi/vim25"
	"github.com/vmware/govmomi/vim25/mo"

	//"github.com/vmware/govmomi/vim25/types"
	"math"

	"github.com/vmware/govmomi/property"
	"github.com/vmware/govmomi/view"
)

func isKeyInList(intval int32, list []int32) bool {
	for _, v := range list {
		if v == intval {
			return true

		}
	}
	return false
}

const MB = 1048576
const GB = 1073741824
const TB = 1099511627776

/*
	Specific VsphereDump Code
*/

func VsphereDump(ctx context.Context, c *vim25.Client, fileName string, shadow bool, defaultsfile string) error {
	defaultWorkLoad := NewWorkLoad()
	if defaultsfile != "" {
		b, derr := ioutil.ReadFile(defaultsfile)
		if derr == nil {
			vseformat := VSEFormat{}
			jerr := json.Unmarshal(b, &vseformat)
			if jerr == nil && len(vseformat.Workloads) > 0 {
				defaultWorkLoad = vseformat.Workloads[0]
			}
		}
	}

	shadowGroups := []ShadowGroup{}
	workloads := []Workload{}

	// Create view of VirtualMachine objects
	m := view.NewManager(c)

	v, err := m.CreateContainerView(ctx, c.ServiceContent.RootFolder, []string{"VirtualMachine", "ComputeResource"}, true)
	if err != nil {
		return err
	}

	defer v.Destroy(ctx)

	/*
		// Retrieve summary property for all machines
		// Reference: http://pubs.vmware.com/vsphere-60/topic/com.vmware.wssdk.apiref.doc/vim.VirtualMachine.html
		var vms []mo.VirtualMachine
		err = v.Retrieve(ctx, []string{"VirtualMachine"}, nil, &vms)
		if err != nil {
			return err
		}

		// Print summary per vm (see also: govc/vm/info.go)
		// https://github.com/vmware/govmomi/blob/master/vim25/mo/mo.go
		// https://vdc-download.vmware.com/vmwb-repository/dcr-public/b50dcbbf-051d-4204-a3e7-e1b618c1e384/538cf2ec-b34f-4bae-a332-3820ef9e7773/vim.vm.GuestInfo.html

	*/

	pc := property.DefaultCollector(c)

	var computeresources []mo.ComputeResource
	err = v.Retrieve(ctx, []string{"ComputeResource"}, nil, &computeresources)
	if err != nil {
		return err
	}

	for _, cluster := range computeresources {
		shadowVMs := []ShadowVM{}

		var rootrp mo.ResourcePool
		err = pc.RetrieveOne(ctx, *cluster.ResourcePool, nil, &rootrp)

		if err == nil {
			workload := defaultWorkLoad
			workload.WorkLoadName = cluster.Name

			var vms []mo.VirtualMachine
			err = pc.Retrieve(ctx, rootrp.Vm, nil, &vms)
			if err == nil {
				for _, vm := range vms {
					isWin := (len(vm.Config.GuestId) > 3 && vm.Config.GuestId[0:3] == "win")

					var vmdiskalloc int64 = 0
					var vmpartsize int64 = 0

					//Counts also vswap etc
					//fmt.Printf("%s\n",vm.Config.Name)
					//for _, ds := range vm.Storage.PerDatastoreUsage {
					//	 vmalloc += ds.Committed
					//}

					workload.VMVMDKQuantity += len(vm.LayoutEx.Disk)

					disks := []ShadowDisk{}
					parts := []ShadowPart{}

					if len(vm.Guest.Disk) > 0 {
						for _, d := range vm.Guest.Disk {
							used := (d.Capacity - d.FreeSpace)
							vmpartsize += used
							if shadow {
								parts = append(parts, ShadowPart{d.DiskPath, d.FreeSpace, d.Capacity, used})
							}
						}
					}

					//get a list of disk
					diskmap := []int32{}
					for _, d := range vm.LayoutEx.Disk {
						//fmt.Println(d)
						for _, c := range d.Chain {
							for _, fk := range c.FileKey {
								diskmap = append(diskmap, fk)
							}
						}
					}

					for _, f := range vm.LayoutEx.File {
						if isKeyInList(f.Key, diskmap) {
							vmdiskalloc += f.Size
							if shadow {
								disks = append(disks, ShadowDisk{f.Name, f.Size})
							}
						}
					}

					var shadowAssigned ShadowAssigned
					if vmpartsize == 0 {
						workload.WorkLoadCapBytes += vmdiskalloc
						workload.ProcessCapBytes += vmdiskalloc
						shadowAssigned = ShadowAssigned{vmdiskalloc, vmdiskalloc}
					} else if isWin {
						workload.WorkLoadCapBytes += vmpartsize
						workload.ProcessCapBytes += vmpartsize
						shadowAssigned = ShadowAssigned{vmpartsize, vmpartsize}
					} else {
						workload.WorkLoadCapBytes += vmpartsize
						workload.ProcessCapBytes += vmdiskalloc
						shadowAssigned = ShadowAssigned{vmpartsize, vmdiskalloc}
					}
					fmt.Printf("Rapporting (Win %t) VM %20s %10d %10d\n", isWin, vm.Config.Name, vmdiskalloc/GB, vmpartsize/GB)
					workload.WorkLoadCapBytes += vmdiskalloc

					if shadow {
						shadowVMs = append(shadowVMs, ShadowVM{vm.Config.Name, vm.Config.GuestId, isWin, disks, parts, vmdiskalloc / GB, vmpartsize / GB, shadowAssigned, cluster.Name})
					}
				}
				workload.VMQuantity = len(vms)
				workload.VMVMDKRatio = int(math.Ceil(float64(workload.VMVMDKQuantity) / float64(workload.VMQuantity)))

				fmt.Printf("VMs : %10d\t\tDisks : %10d\t\tD/V : %10d", workload.VMQuantity, workload.VMVMDKQuantity, workload.VMVMDKRatio)

				workload.WorkLoadCap = int64(math.Ceil(float64(workload.WorkLoadCapBytes) / TB))
				workload.ProcessCap = int64(math.Ceil(float64(workload.ProcessCapBytes) / TB))

				if shadow {
					shadowGroups = append(shadowGroups, ShadowGroup{shadowVMs, cluster.Name})
				}

				workloads = append(workloads, workload)
			} else {
				return err
			}
		} else {
			return err
		}
	}

	vseformat := VSEFormat{workloads}

	txt, _ := json.Marshal(vseformat)

	ioutil.WriteFile(fileName, txt, 0777)

	if shadow {
		txt, _ := json.Marshal(ShadowFile{shadowGroups})
		ioutil.WriteFile(fileName+".shadow.json", txt, 0777)
	}

	return nil
}
