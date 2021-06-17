package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"strconv"
	"io/ioutil"
//	"net/http"
	"time"
	"github.com/OpenNebula/one/src/oca/go/src/goca"
)

func NewConfigw(user string, password string, endpoint string) goca.OneConfig {
	var authToken string
	var oneAuthPath string

	oneXmlrpc := endpoint

	if user == "" && password == "" {
		oneAuthPath = os.Getenv("ONE_AUTH")
		if oneAuthPath == "" {
			oneAuthPath = os.Getenv("HOME") + "/.one/one_auth"
		}

		token, err := ioutil.ReadFile(oneAuthPath)
		if err == nil {
			authToken = strings.TrimSpace(string(token))
		} else if authToken == "" {
			token2, err := ioutil.ReadFile("/var/lib/one/.one/one_auth")
			if err == nil {
				authToken = strings.TrimSpace(string(token2))
			}
		}else{
			authToken = ""
		}
	} else {
		authToken = user + ":" + password
	}

	if oneXmlrpc == "" {
		oneXmlrpc = os.Getenv("ONE_XMLRPC")
		if oneXmlrpc == "" {
			oneXmlrpc = "http://localhost:2633/RPC2"
		}
	}

	config := goca.OneConfig{
		Token:    authToken,
		Endpoint: oneXmlrpc,
	}

	return config
}
func main() {
	client := goca.NewDefaultClient(
		NewConfigw("","",""),
	)
	controller := goca.NewController(client)

	// Get short informations of the VMs
	vms, err := controller.VMs().Info()
	if err != nil {
		log.Fatal(err)
	}
	for _, vm := range vms.VMs {
		test,test2, err := vm.State()
		if err != nil {
			log.Fatal(err)
		}
		vm2,_ := controller.VM(vm.ID).Info(false)
	//	fmt.Printf("VM: %s",vm)
		stime := time.Unix(int64(vm.STime),0)
		fmt.Printf("ID: %s, Name: %s, State: %s, LCM_STATE: %s STIME: %s ETIME: %s UPTIME: %s\n", strconv.Itoa(vm.ID), vm.Name, test, test2, stime, vm.ETime,time.Since(stime))
		nics := vm2.Template.GetNICs()
		disks := vm.Template.GetDisks()
		cpu,_ := vm.Template.GetCPU()
		vcpu,_  := vm.Template.GetVCPU()
		memory,_ := vm.Template.GetMemory()
		os_arch,_ := vm2.Template.GetOS("ARCH")
		os_boot,_ := vm2.Template.GetOS("BOOT")
		os_machine,_ := vm2.Template.GetOS("MACHINE")
		graphic_port,_ := vm.Template.GetIOGraphic("PORT")
		graphic_type,_ := vm.Template.GetIOGraphic("TYPE")
		for _, nic := range nics {
			//fmt.Printf("%s",nic)
			m,_ := nic.Get("MAC")
			i,_ := nic.Get("IP")
			n,_ := nic.Get("NETWORK")
			p,_ := nic.Get("PHYDEV")
			mo,_ := nic.Get("MODEL")
			t,_ := nic.Get("TARGET")
			b,_ := nic.Get("BRIDGE")
			bt,_ := nic.Get("BRIDGE_TYPE")
			v,_ := nic.Get("VLAN_ID")
			fmt.Printf("NETWORK: %s MAC: %s IP: %s PHYDEV: %s MODEL: %s TARGET: %s BRIDGE: %s BRIDGE_TYPE: %s VLAN_ID: %s\n",n,m,i,p,mo,t,b,bt,v)
		}
		for _, disk := range disks {
		//	fmt.Printf("%s\n",disk)
			disk_id,_ := disk.Get("DISK_ID")
			image_id,_ := disk.Get("IMAGE_ID")
			image,_ := disk.Get("IMAGE")
			target,_  := disk.Get("TARGET")
			size,_ := disk.Get("SIZE")
			fmt.Printf("DISK_ID: %s, IMAGE_ID: %s, IMAGE: %s, TARGET: %s, SIZE: %s\n", disk_id,image_id,image,target,size)
		}
		fmt.Printf("GRAPHIC_TYPE: %s GRAPHIC_PORT: %s\n",graphic_type,graphic_port)
		fmt.Printf("OS_ARCH: %s, OS_BOOT: %s, OS_MACHINE: %s\n",os_arch,os_boot,os_machine)
		fmt.Printf("CPU: %s\n",fmt.Sprintf("%.0f",cpu))
		fmt.Printf("VCPU: %s\n",strconv.Itoa(vcpu))
		fmt.Printf("MEMORY: %s\n",memory)
		records := vm.HistoryRecords
		for _, rec := range records {
			fmt.Printf("HOST: %s START-TIME: %s END-TIME: %s\n\n",rec.Hostname,rec.RSTime,rec.RETime)
		}
		mons := vm2.MonitoringInfos
		fmt.Printf("MONITORING INFO: %s\n",mons)
		mons2 := vm2.MonitoringInfos.GetVectors("DISK_SIZE")
		fmt.Printf("MONITORING PAIR: %s\n",mons2)
		for _, vec := range mons2 {
			fmt.Printf("vectors: %s\n",vec.GetStrs("SIZE")[0])
		}
		//fmt.Printf("MONITORING PAIR ERROR: %s\n",err)
		//fmt.Printf("%s\n",vm.Template.GetNICs())
	}
}
