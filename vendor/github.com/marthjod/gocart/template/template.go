package template

import (
	"encoding/xml"
	"fmt"

	"github.com/marthjod/gocart/api"
	"github.com/marthjod/gocart/image"
	"github.com/marthjod/gocart/vnet"
)

// Template represents a VM template.
type Template struct {
	ID       int          `xml:"ID"`
	Name     string       `xml:"NAME"`
	Uname    string       `xml:"UNAME"`
	RegTime  int          `xml:"REGTIME"`
	Template HostTemplate `xml:"TEMPLATE"`
	Memory   int          `xml:"MEMORY"`
	VMID     int          `xml:"VMID"`
	Disk     []image.Disk `xml:"DISK"`
	CPU      string       `xml:"CPU"`
}

// HostTemplate represents a host template.
type HostTemplate struct {
	XMLName        xml.Name     `xml:"TEMPLATE"`
	CPU            string       `xml:"CPU"`
	Disk           []image.Disk `xml:"DISK"`
	Memory         string       `xml:"MEMORY"`
	Name           string       `xml:"NAME"`
	Nics           []vnet.NIC   `xml:"NIC"`
	VCPU           string       `xml:"VCPU"`
	Datacenter     string       `xml:"DATACENTER"`
	Requirements   string       `xml:"REQUIREMENTS"`
	DSRequirements string       `xml:"SCHED_DS_REQUIREMENTS"`
	Items          Tags         `xml:",any"`
}

// UserTemplate represents a user template.
type UserTemplate struct {
	Items Tags `xml:",any"`
}

// Tag is an XML tag.
type Tag struct {
	XMLName xml.Name
	Content string `xml:",chardata"`
}

// Tags is a list of Tags.
type Tags []Tag

// GetCustom returns values from custom-defined XML tags.
func (tags Tags) GetCustom(tagName string) (string, error) {
	for _, tag := range tags {
		if tagName == tag.XMLName.Local {
			return tag.Content, nil
		}
	}
	return "", fmt.Errorf("tag %s not found", tagName)
}

// Instantiate http://docs.opennebula.org/4.12/integration/system_interfaces/api.html#one-template-instantiate
func (t *Template) Instantiate(c *api.RPC, name string) error {
	return c.Call(t, "one.template.instantiate", []interface{}{c.AuthString, t.ID, name, false, ""})
}

func (t Template) String() string {
	return t.Name
}
