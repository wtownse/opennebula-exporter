package main

import (
        "fmt"
	"github.com/OpenNebula/one/src/oca/go/src/goca"
	"github.com/OpenNebula/one/src/oca/go/src/goca/schemas/vm"
  	"net/http"
  	"gopkg.in/alecthomas/kingpin.v2"
  	"os"
  	log "github.com/Sirupsen/logrus"
  	"github.com/prometheus/client_golang/prometheus"
  	"github.com/prometheus/client_golang/prometheus/promhttp"
	"strconv"
)

const (
       endpoint = "http://localhost:2633/RPC2"
        user     = "oneadmin"
       pass     = "opennebula"
)

type PoweroneNode struct {
        Id              int
        State           string
        Name            string
        Url             string
        LogIndex        int
        Commit          int
	Term		int
}

type PoweroneCluster struct {
        PoweroneNodes   []PoweroneNode
}

type oneCollector struct {
        oneServerRaftStatusTermDesc		*prometheus.Desc
        oneServerRaftStatusVotedforDesc		*prometheus.Desc
        oneServerRaftStatusCommitDesc		*prometheus.Desc
        oneServerRaftStatusLogIndexDesc		*prometheus.Desc
	oneServerRaftStatusFedlogIndexDesc	*prometheus.Desc
	oneServerRaftStatusSyncedDesc		*prometheus.Desc
	oneVmListDesc				*prometheus.Desc
	oneVmNicListDesc			*prometheus.Desc
	oneVmDiskListDesc			*prometheus.Desc
	oneVmGraphicDesc			*prometheus.Desc
	oneVmOsDesc				*prometheus.Desc
}
		
func newoneCollector() *oneCollector  {
        return &oneCollector   {
		oneServerRaftStatusTermDesc: prometheus.NewDesc("oneServerRaftStatusTerm",
                        "Shows opennebula raft term",
                        []string{"metricname","host","state","id","endpoint"},
                        nil,
                ),
		oneServerRaftStatusVotedforDesc: prometheus.NewDesc("oneServerRaftStatusVotedfor",
                        "Shows opennebula raft voted for",
                        []string{"metricname","host","state","id","endpoint"},
                        nil,
                ),
		oneServerRaftStatusCommitDesc: prometheus.NewDesc("oneServerRaftStatusCommit",
                        "Shows opennebula raft commit",
                        []string{"metricname","host","state","id","endpoint"},
                        nil,
                ),
		oneServerRaftStatusLogIndexDesc: prometheus.NewDesc("oneServerRaftStatusLogIndex",
                        "Shows opennebula raft log index",
                        []string{"metricname","host","state","id","endpoint"},
                        nil,
                ),
		oneServerRaftStatusFedlogIndexDesc: prometheus.NewDesc("oneServerRaftStatusFedlogIndex",
                        "Shows opennebula raft fedlog index",
                        []string{"metricname","host","state","id","endpoint"},
                        nil,
                ),
		oneServerRaftStatusSyncedDesc: prometheus.NewDesc("oneServerRaftStatusSyncedDesc",
                        "Shows opennebula raft fedlog index",
                        []string{"metricname","host","state","id","endpoint","term","index","commit"},
                        nil,
                ),
                oneVmListDesc: prometheus.NewDesc("oneVmListDesc",
                        "Shows opennebula virtual machine status",
                        []string{"metricname","id","name","state","lcm_state","host","cpu","vcpu","memory"},
                        nil,
                ),
                oneVmNicListDesc: prometheus.NewDesc("oneVmNicListDesc",
                        "Shows opennebula virtual machine interface(s) info",
                        []string{"metricname","id","name","mac","network","ip","phydev","model","target","bridge","bridge_type","vlan_id"},
                        nil,
                ),
                oneVmDiskListDesc: prometheus.NewDesc("oneVmDiskListDesc",
                        "Shows opennebula virtual machine disk(s) info",
                        []string{"metricname","id","name","disk_id","image_id","image","target","size","space_used"},
                        nil,
                ),
                oneVmGraphicDesc: prometheus.NewDesc("oneVmGraphicDesc",
                        "Shows opennebula virtual machine graphic info",
                        []string{"metricname","id","name","graphic_type","hostname","graphic_port"},
                        nil,
                ),
                oneVmOsDesc: prometheus.NewDesc("oneVmOsDesc",
                        "Shows opennebula virtual machine os info",
                        []string{"metricname","id","name","os_arch","os_boot","os_machine"},
                        nil,
                ),
        }
}

func (collector *oneCollector) Describe(ch chan<- *prometheus.Desc) {
        ch <- collector.oneServerRaftStatusTermDesc
        ch <- collector.oneServerRaftStatusVotedforDesc
        ch <- collector.oneServerRaftStatusCommitDesc
        ch <- collector.oneServerRaftStatusLogIndexDesc
        ch <- collector.oneServerRaftStatusFedlogIndexDesc
	ch <- collector.oneServerRaftStatusSyncedDesc
	ch <- collector.oneVmListDesc
	ch <- collector.oneVmNicListDesc
	ch <- collector.oneVmDiskListDesc
	ch <- collector.oneVmGraphicDesc
	ch <- collector.oneVmOsDesc
}


func (collector *oneCollector) Collect(ch chan<- prometheus.Metric) {
	conf := goca.NewConfig(user,pass,endpoint)
	client := goca.NewDefaultClient(conf)
	id := 0
	controller := goca.NewController(client)

        // Retrieve zone info
        zone, err := controller.Zone(id).Info(false)
	if err != nil {
		log.Print("Zone id %d: %s", id, err)
	}
	var pclus PoweroneCluster
	for _, server := range zone.ServerPool {
		conf.Endpoint = server.Endpoint
		client = goca.NewDefaultClient(conf)
		controller.Client = client
		// Fetch the raft status of the server behind the endpoint
		status, err := controller.Zones().ServerRaftStatus()
		if err != nil {
			log.Print("Server raft status endpoint %s: %s", server.Endpoint, err)
                	var PoneNode PoweroneNode
                	PoneNode.Id = server.ID
                	PoneNode.State = "UNREACHABLE"
                	PoneNode.Name = server.Name
                	PoneNode.Url = server.Endpoint
                	PoneNode.LogIndex = -1
                	PoneNode.Commit = -1
                	pclus.PoweroneNodes = append(pclus.PoweroneNodes,PoneNode)
			continue
		}

		// Display the Raft state of the server: Leader, Follower, Candidate, Error
		state, err := status.State()
		if err != nil {
			log.Print("Server raft state %d: %s", status.StateRaw, err)
		}
                var PoneNode PoweroneNode
                PoneNode.Id = server.ID
                PoneNode.State = state.String()
                PoneNode.Name = server.Name
                PoneNode.Url = server.Endpoint
                PoneNode.LogIndex = status.LogIndex
                PoneNode.Commit = status.Commit
                pclus.PoweroneNodes = append(pclus.PoweroneNodes,PoneNode)

        ch <- prometheus.MustNewConstMetric(
                collector.oneServerRaftStatusTermDesc,
                prometheus.CounterValue,
                float64(status.Term),
                "Term",
                server.Name,
		state.String(),
                strconv.Itoa(server.ID),
		server.Endpoint)
        ch <- prometheus.MustNewConstMetric(
                collector.oneServerRaftStatusVotedforDesc,
                prometheus.CounterValue,
                float64(status.Votedfor),
                "Vote",
                server.Name,
		state.String(),
                strconv.Itoa(server.ID),
		server.Endpoint)
        ch <- prometheus.MustNewConstMetric(
                collector.oneServerRaftStatusCommitDesc,
                prometheus.CounterValue,
                float64(status.Commit),
                "Commit",
                server.Name,
		state.String(),
                strconv.Itoa(server.ID),
		server.Endpoint)
        ch <- prometheus.MustNewConstMetric(
                collector.oneServerRaftStatusLogIndexDesc,
                prometheus.CounterValue,
                float64(status.LogIndex),
                "LogIndex",
                server.Name,
		state.String(),
                strconv.Itoa(server.ID),
		server.Endpoint)
        ch <- prometheus.MustNewConstMetric(
                collector.oneServerRaftStatusFedlogIndexDesc,
                prometheus.CounterValue,
                float64(status.FedlogIndex),
                "FedlogIndex",
                server.Name,
		state.String(),
                strconv.Itoa(server.ID),
		server.Endpoint)
	}
	
        for _, node := range pclus.PoweroneNodes {
                if node.State == "LEADER"{
			synced := 0
			MIN := node.LogIndex - 10
                        for _, ncompare := range pclus.PoweroneNodes {
                                if node.State == ncompare.State {
                                        synced = 1
                                } else if ncompare.State == "FOLLOWER" {
                                        if ncompare.LogIndex >= MIN {
                                                synced = 1
                                        } else {
                                                synced = 0
                                        }
                                } else {
                                        synced = -1
                                }

        ch <- prometheus.MustNewConstMetric(
                collector.oneServerRaftStatusSyncedDesc,
                prometheus.CounterValue,
                float64(synced),
                "Synced",
                ncompare.Name,
                ncompare.State,
                strconv.Itoa(ncompare.Id),
                node.Url,
		strconv.Itoa(ncompare.Term),
		strconv.Itoa(ncompare.LogIndex),
		strconv.Itoa(ncompare.Commit))
                        }
                }
        }
        // Get short informations of the VMs
        vms, err := controller.VMs().Info()
        if err != nil {
                log.Fatal(err)
        }
        for _, onevm := range vms.VMs {
                state,lcm_state, err := onevm.StateString()
                if err != nil {
                        log.Fatal(err)
                }
                vm2,_ := controller.VM(onevm.ID).Info(false)
                nics := vm2.Template.GetNICs()
                disks := vm2.Template.GetDisks()
                cpu,_ := onevm.Template.GetCPU()
                vcpu,_  := onevm.Template.GetVCPU()
                memory,_ := onevm.Template.GetMemory()
                os_arch,err := vm2.Template.GetOS("ARCH")
                os_boot,_ := vm2.Template.GetOS("BOOT")
                os_machine,_ := vm2.Template.GetOS("MACHINE")
                graphic_port,_ := onevm.Template.GetIOGraphic("PORT")
                graphic_type,_ := onevm.Template.GetIOGraphic("TYPE")
		disk_used := vm2.MonitoringInfos.GetVectors("DISK_SIZE")
		records := onevm.HistoryRecords
		var rhost vm.HistoryRecord
                for _, rec := range records {
                        rhost = rec
                }
		var running int = 0
		if lcm_state == "RUNNING" {
			running = 1
		}
	name := onevm.Name
        ch <- prometheus.MustNewConstMetric(
                collector.oneVmListDesc,
                prometheus.CounterValue,
                float64(running),
                "oneVMList",
		strconv.Itoa(onevm.ID),
                name,
                state,
                lcm_state,
                rhost.Hostname,
		fmt.Sprintf("%.0f",cpu),
		strconv.Itoa(vcpu),
		strconv.Itoa(memory))

                for _, nic := range nics {
                        mac,_ := nic.Get("MAC")
                        ip,_ := nic.Get("IP")
                        net,_ := nic.Get("NETWORK")
                        phy,_ := nic.Get("PHYDEV")
                        model,_ := nic.Get("MODEL")
                        target,_ := nic.Get("TARGET")
                        bridge,_ := nic.Get("BRIDGE")
                        bridge_type,_ := nic.Get("BRIDGE_TYPE")
                        vlan,_ := nic.Get("VLAN_ID")
        		ch <- prometheus.MustNewConstMetric(
                		collector.oneVmNicListDesc,
                		prometheus.CounterValue,
                		float64(1),
                		"oneVmNicList",
				strconv.Itoa(onevm.ID),
                		onevm.Name,
                		mac,
				net,
				ip,
				phy,
				model,
				target,
				bridge,
				bridge_type,
				vlan)
                }
                for _, disk := range disks {
                        disk_id,_ := disk.Get("DISK_ID")
                        image_id,_ := disk.Get("IMAGE_ID")
                        image,_ := disk.Get("IMAGE")
                        target,_  := disk.Get("TARGET")
                        size,_ := disk.Get("SIZE")
                        //fmt.Print(disk)
                        //fmt.Printf("DISK_ID: %s, IMAGE_ID: %s, IMAGE: %s, TARGET: %s, SIZE: %s\n", disk_id,image_id,image,target,size)
        		s,_ := strconv.ParseFloat(size,64)
			used_space := "0"
			for _, d := range disk_used {
				if d.GetStrs("ID")[0] == disk_id {
				used_space = d.GetStrs("SIZE")[0]
				}
			}
			ch <- prometheus.MustNewConstMetric(
                		collector.oneVmDiskListDesc,
                		prometheus.CounterValue,
                		float64(s),
                		"oneVmDiskListDesc",
				strconv.Itoa(onevm.ID),
                		onevm.Name,
				disk_id,
				image_id,
				image,
				target,
				size,
				used_space)
                }
		ch <- prometheus.MustNewConstMetric(
                	collector.oneVmGraphicDesc,
                	prometheus.CounterValue,
                	float64(1),
                	"oneVmGraphic",
			strconv.Itoa(onevm.ID),
                	onevm.Name,
			graphic_type,
			rhost.Hostname,
			graphic_port)
		ch <- prometheus.MustNewConstMetric(
                	collector.oneVmOsDesc,
                	prometheus.CounterValue,
                	float64(1),
                	"oneVmOs",
			strconv.Itoa(onevm.ID),
                	onevm.Name,
			os_arch,
			os_boot,
			os_machine)
        }

}

func main() {
var (
app           = kingpin.New("opennebula_exporter", "Prometheus metrics exporter for opennebula")
listenAddress = app.Flag("web.listen-address", "Address to listen on for web interface and telemetry.").Default(":9701").String()
metricsPath   = app.Flag("web.telemetry-path", "Path under which to expose metrics.").Default("/metrics").String()
)
kingpin.MustParse(app.Parse(os.Args[1:]))
        onecollector := newoneCollector()
        
        prometheus.MustRegister(onecollector)

http.Handle(*metricsPath, promhttp.Handler())
http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte(`
                <html>
                <head><title>Opennebula Exporter</title></head>
                <body>
                <h1>Opennebula Exporter</h1>
                <p><a href='` + *metricsPath + `'>Metrics</a></p>
                </body>
                </html>`))
        })
	log.Info("Beginning to serve on port ",*listenAddress)
	log.Fatal(http.ListenAndServe(*listenAddress, nil))
}
