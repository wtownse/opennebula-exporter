package hostpool

import (
	"encoding/xml"
	"io"

	"github.com/marthjod/gocart/api"
	"github.com/marthjod/gocart/host"
	"github.com/marthjod/gocart/vmpool"
)

// HostPool represents a host pool.
type HostPool struct {
	XMLName xml.Name     `xml:"HOST_POOL"`
	Hosts   []*host.Host `xml:"HOST"`
}

// Info http://docs.opennebula.org/4.12/integration/system_interfaces/api.html#one-hostpool-info
func (p *HostPool) Info(c *api.RPC) error {
	return c.Call(p, "one.hostpool.info", []interface{}{c.AuthString})
}

// MapVMs ...
func (p *HostPool) MapVMs(vmpool *vmpool.VMPool) {
	for _, host := range p.Hosts {
		host.MapVMs(vmpool)
	}
}

// FromReader reads into a host pool.
func FromReader(r io.Reader) (*HostPool, error) {
	pool := HostPool{}
	dec := xml.NewDecoder(r)
	if err := dec.Decode(&pool); err != nil {
		return nil, err
	}
	return &pool, nil
}

// GetHostsInCluster returns a pool of hosts in the provided cluster.
func (p *HostPool) GetHostsInCluster(cluster string) *HostPool {
	var (
		hostpool HostPool
	)
	for _, host := range p.Hosts {
		if host.Cluster == cluster {
			hostpool.Hosts = append(hostpool.Hosts, host)
		}
	}
	return &hostpool
}

// FilterHostsByStates returns host pool containing only hosts in one of the provided states.
func (p *HostPool) FilterHostsByStates(states ...host.State) *HostPool {
	var (
		hp HostPool
	)
	for _, host := range p.Hosts {
		for _, state := range states {
			if host.State == state {
				hp.Hosts = append(hp.Hosts, host)
				continue
			}
		}
	}
	return &hp
}

// FilterOutEmptyHosts filters out hosts without VMs.
func (p *HostPool) FilterOutEmptyHosts() *HostPool {
	var (
		hp HostPool
	)
	for _, host := range p.Hosts {
		if !host.IsEmpty() {
			hp.Hosts = append(hp.Hosts, host)
		}
	}
	return &hp
}
