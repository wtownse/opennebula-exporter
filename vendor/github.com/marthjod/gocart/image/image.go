package image

import "encoding/xml"

// Disk represents a disk.
type Disk struct {
	XMLName xml.Name `xml:"DISK"`
	Name    string   `xml:"IMAGE"`
	ID      int      `xml:"IMAGE_ID"`
}

// Image represents an image.
type Image struct {
	XMLName     xml.Name `xml:"IMAGE"`
	ID          int      `xml:"ID"`
	Name        string   `xml:"NAME"`
	Datastore   string   `xml:"DATASTORE"`
	DatastoreID int      `xml:"DATASTORE_ID"`
	RunningVMs  int      `xml:"RUNNING_VMS"`
}
