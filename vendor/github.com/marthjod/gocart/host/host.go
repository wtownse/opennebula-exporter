package host

import (
	"encoding/xml"

	"github.com/marthjod/gocart/template"
	"github.com/marthjod/gocart/vmpool"
)

// State represents host state
type State int

//go:generate stringer -type=State
const (
	Init                State = iota
	MonitoringMonitored State = iota // Currently monitoring, previously MONITORED
	Monitored           State = iota
	Error               State = iota
	Disabled            State = iota
	MonitoringError     State = iota // Currently monitoring, previously ERROR
	MonitoringInit      State = iota // Currently monitoring, previously initialized
	MonitoringDisabled  State = iota // Currently monitoring, previously DISABLED
)

// Host represents an OpenNebula host.
type Host struct {
	XMLName   xml.Name              `xml:"HOST"`
	ID        int                   `xml:"ID"`
	Name      string                `xml:"NAME"`
	State     State                 `xml:"STATE"`
	Cluster   string                `xml:"CLUSTER"`
	ClusterID int                   `xml:"CLUSTER_ID"`
	Template  template.HostTemplate `xml:"TEMPLATE"`
	VMIDs     []int                 `xml:"VMS>ID"`
	VMPool    *vmpool.VMPool
}

// MapVMs populates the host's VM pool with its VMs.
func (h *Host) MapVMs(vmpool *vmpool.VMPool) {
	h.VMPool = vmpool.GetVMsByID(h.VMIDs...)
}

// IsEmpty checks if a host has no VMs.
func (h *Host) IsEmpty() bool {
	return len(h.VMIDs) == 0
}

// String returns a host's short strings representation.
func (h *Host) String() string {
	return h.Name
}
