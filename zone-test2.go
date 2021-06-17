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

type PoweroneNode struct {
	Id		int
	State		string
	Name		string
	Url		string
	LogIndex	int
	Commit		int
}

type PoweroneCluster struct {
	PoweroneNodes	[]PoweroneNode
}

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
	var test PoweroneCluster
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
		var PoneNode PoweroneNode
		PoneNode.Id = server.ID
		PoneNode.State = state.String()
		PoneNode.Name = server.Name
		PoneNode.Url = server.Endpoint
		PoneNode.LogIndex = status.LogIndex
		PoneNode.Commit = status.Commit
		test.PoweroneNodes = append(test.PoweroneNodes,PoneNode)
	}
	for _, node := range test.PoweroneNodes {
		if node.State == "LEADER"{
			fmt.Printf("data: %s\n",node)
			for _, ncompare := range test.PoweroneNodes {
				if node.State == ncompare.State {
					fmt.Printf("We are Groot: %s\n",ncompare)
				} else if ncompare.State == "FOLLOWER" {
					if ncompare.LogIndex == node.LogIndex {
						fmt.Printf("FOLLOWER IN SYNC: %s\n",ncompare)
					} else {
						fmt.Printf("FOLLOWER OUT OF SYNC: %s\n",ncompare)
					}
				} else {
					fmt.Printf("NOT IN LEADER OR FOLLOWER ROLE: %s\n",ncompare)
				}
			}
		}
	}

}
