package main

import (
        "fmt"
//	"github.com/OpenNebula/goca"
	"github.com/OpenNebula/one/src/oca/go/src/goca"
        "log"
	"strconv"
)

const (
       endpoint = "http://localhost:2633/RPC2"
        user     = "oneadmin"
       pass     = "opennebula"
)

func main() {
	conf := goca.NewConfig(user,pass,endpoint)
	client := goca.NewDefaultClient(conf)
	id := 0
	controller := goca.NewController(client)

        // Retrieve zone info
        zone, err := controller.Zone(id).Info(false)
	if err != nil {
		log.Fatalf("Zone id %d: %s", id, err)
	}
	for _, server := range zone.ServerPool {
		fmt.Println(server)
		var HOST,_ = controller.Host(server.ID).Info(false)
		fmt.Println(HOST.Name)
		conf.Endpoint = server.Endpoint
		client = goca.NewDefaultClient(conf)
		controller.Client = client
		// Fetch the raft status of the server behind the endpoint
		status, err := controller.Zones().ServerRaftStatus()
		if err != nil {
			log.Fatalf("Server raft status endpoint %s: %s", server.Endpoint, err)
		}

		// Display the Raft state of the server: Leader, Follower, Candidate, Error
		state, err := status.State()
		if err != nil {
			log.Fatalf("Server raft state %d: %s", status.StateRaw, err)
		}
		fmt.Printf("server: %s, state: %s commit_index: %s log_index: %s\n", server.Name, state.String(),strconv.Itoa(status.Commit),strconv.Itoa(status.LogIndex))
		fmt.Printf(strconv.Itoa(HOST.Share.FreeDisk))
	}
	// Get short informations of the VMs
	vms, err := controller.VMs().Info()
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < len(vms.VMs); i++ {
		// This Info method, per VM instance, give us detailed informations on the instance
		// Check xsd files to see the difference
		vm, err := controller.VM(vms.VMs[i].ID).Info(false)
		if err != nil {
			log.Fatal(err)
		}

		//Do some others stuffs on vm
		fmt.Printf("%+v\n", vm.HistoryRecords)
	}
}
