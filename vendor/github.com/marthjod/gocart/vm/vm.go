package vm

import (
	"encoding/xml"
	"io"

	"github.com/marthjod/gocart/image"
	"github.com/marthjod/gocart/template"
)

// State represents VM state
type State int

// LCMState represents LCM (lifecycle manager) state
type LCMState int

//go:generate stringer -type=State
const (
	Init       State = iota
	Pending    State = iota
	Hold       State = iota
	Active     State = iota
	Stopped    State = iota
	Suspended  State = iota
	Done       State = iota
	Failed     State = iota
	Poweroff   State = iota
	Undeployed State = iota
)

//go:generate stringer -type=LCMState
const (
	LcmInit                      LCMState = iota
	Prolog                       LCMState = iota
	Boot                         LCMState = iota
	Running                      LCMState = iota
	Migrate                      LCMState = iota
	SaveStop                     LCMState = iota
	SaveSuspend                  LCMState = iota
	SaveMigrate                  LCMState = iota
	PrologMigrate                LCMState = iota
	PrologResume                 LCMState = iota
	EpilogStop                   LCMState = iota
	Epilog                       LCMState = iota
	Shutdown                     LCMState = iota
	CleanupResubmit              LCMState = iota
	Unknown                      LCMState = iota
	Hotplug                      LCMState = iota
	ShutdownPoweroff             LCMState = iota
	BootUnknown                  LCMState = iota
	BootPoweroff                 LCMState = iota
	BootSuspended                LCMState = iota
	BootStopped                  LCMState = iota
	CleanupDelete                LCMState = iota
	HotplugSnapshot              LCMState = iota
	HotplugNic                   LCMState = iota
	HotplugSaveas                LCMState = iota
	HotplugSaveasPoweroff        LCMState = iota
	HotplugSaveasSuspended       LCMState = iota
	ShutdownUndeploy             LCMState = iota
	EpilogUndeploy               LCMState = iota
	PrologUndeploy               LCMState = iota
	BootUndeploy                 LCMState = iota
	HotplugPrologPoweroff        LCMState = iota
	HotplugEpilogPoweroff        LCMState = iota
	BootMigrate                  LCMState = iota
	BootFailure                  LCMState = iota
	BootMigrateFailure           LCMState = iota
	PrologMigrateFailure         LCMState = iota
	PrologFailure                LCMState = iota
	EpilogFailure                LCMState = iota
	EpilogStopFailure            LCMState = iota
	EpilogUndeployFailure        LCMState = iota
	PrologMigratePoweroff        LCMState = iota
	PrologMigratePoweroffFailure LCMState = iota
	PrologMigrateSuspend         LCMState = iota
	PrologMigrateSuspendFailure  LCMState = iota
	BootUndeployFailure          LCMState = iota
	BootStoppedFailure           LCMState = iota
	PrologResumeFailure          LCMState = iota
	PrologUndeployFailure        LCMState = iota
	DiskSnapshotPoweroff         LCMState = iota
	DiskSnapshotRevertPoweroff   LCMState = iota
	DiskSnapshotDeletePoweroff   LCMState = iota
	DiskSnapshotSuspended        LCMState = iota
	DiskSnapshotRevertSuspended  LCMState = iota
	DiskSnapshotDeleteSuspended  LCMState = iota
	DiskSnapshot                 LCMState = iota
	DiskSnapshotDelete           LCMState = iota
	PrologMigrateUnknown         LCMState = iota
	PrologMigrateUnknownFailure  LCMState = iota
	DiskResize                   LCMState = iota
	DiskResizePoweroff           LCMState = iota
	DiskResizeUndeployed         LCMState = iota
)

// States maps VM state names to their constant State values.
var States = map[string]State{
	"Init":       Init,
	"Pending":    Pending,
	"Hold":       Hold,
	"Active":     Active,
	"Stopped":    Stopped,
	"Suspended":  Suspended,
	"Done":       Done,
	"Failed":     Failed,
	"Poweroff":   Poweroff,
	"Undeployed": Undeployed,
}

// LCMStates maps LCM state names to their constant LCMState values.
var LCMStates = map[string]LCMState{
	"LcmInit":                      LcmInit,
	"Prolog":                       Prolog,
	"Boot":                         Boot,
	"Running":                      Running,
	"Migrate":                      Migrate,
	"SaveStop":                     SaveStop,
	"SaveSuspend":                  SaveSuspend,
	"SaveMigrate":                  SaveMigrate,
	"PrologMigrate":                PrologMigrate,
	"PrologResume":                 PrologResume,
	"EpilogStop":                   EpilogStop,
	"Epilog":                       Epilog,
	"Shutdown":                     Shutdown,
	"CleanupResubmit":              CleanupResubmit,
	"Unknown":                      Unknown,
	"Hotplug":                      Hotplug,
	"ShutdownPoweroff":             ShutdownPoweroff,
	"BootUnknown":                  BootUnknown,
	"BootPoweroff":                 BootPoweroff,
	"BootSuspended":                BootSuspended,
	"BootStopped":                  BootStopped,
	"CleanupDelete":                CleanupDelete,
	"HotplugSnapshot":              HotplugSnapshot,
	"HotplugNic":                   HotplugNic,
	"HotplugSaveas":                HotplugSaveas,
	"HotplugSaveasPoweroff":        HotplugSaveasPoweroff,
	"HotplugSaveasSuspended":       HotplugSaveasSuspended,
	"ShutdownUndeploy":             ShutdownUndeploy,
	"EpilogUndeploy":               EpilogUndeploy,
	"PrologUndeploy":               PrologUndeploy,
	"BootUndeploy":                 BootUndeploy,
	"HotplugPrologPoweroff":        HotplugPrologPoweroff,
	"HotplugEpilogPoweroff":        HotplugEpilogPoweroff,
	"BootMigrate":                  BootMigrate,
	"BootFailure":                  BootFailure,
	"BootMigrateFailure":           BootMigrateFailure,
	"PrologMigrateFailure":         PrologMigrateFailure,
	"PrologFailure":                PrologFailure,
	"EpilogFailure":                EpilogFailure,
	"EpilogStopFailure":            EpilogStopFailure,
	"EpilogUndeployFailure":        EpilogUndeployFailure,
	"PrologMigratePoweroff":        PrologMigratePoweroff,
	"PrologMigratePoweroffFailure": PrologMigratePoweroffFailure,
	"PrologMigrateSuspend":         PrologMigrateSuspend,
	"PrologMigrateSuspendFailure":  PrologMigrateSuspendFailure,
	"BootUndeployFailure":          BootUndeployFailure,
	"BootStoppedFailure":           BootStoppedFailure,
	"PrologResumeFailure":          PrologResumeFailure,
	"PrologUndeployFailure":        PrologUndeployFailure,
	"DiskSnapshotPoweroff":         DiskSnapshotPoweroff,
	"DiskSnapshotRevertPoweroff":   DiskSnapshotRevertPoweroff,
	"DiskSnapshotDeletePoweroff":   DiskSnapshotDeletePoweroff,
	"DiskSnapshotSuspended":        DiskSnapshotSuspended,
	"DiskSnapshotRevertSuspended":  DiskSnapshotRevertSuspended,
	"DiskSnapshotDeleteSuspended":  DiskSnapshotDeleteSuspended,
	"DiskSnapshot":                 DiskSnapshot,
	"DiskSnapshotDelete":           DiskSnapshotDelete,
	"PrologMigrateUnknown":         PrologMigrateUnknown,
	"PrologMigrateUnknownFailure":  PrologMigrateUnknownFailure,
	"DiskResize":                   DiskResize,
	"DiskResizePoweroff":           DiskResizePoweroff,
	"DiskResizeUndeployed":         DiskResizeUndeployed,
}

// GetState returns a VM state for a given string.
func GetState(state string) State {
	return States[state]
}

// GetLCMState returns an LCMState for a given string.
func GetLCMState(state string) LCMState {
	return LCMStates[state]
}

// VM ...
// Node (current OpenNebula node the VM is running on) determination is best effort:
// in case of multiple hosts, this *might* not be 100% reliable, ie. pick the current one,
// although tests reproducibly showed correct results.
type VM struct {
	XMLName      xml.Name              `xml:"VM"`
	ID           int                   `xml:"ID"`
	Name         string                `xml:"NAME"`
	CPU          int                   `xml:"CPU"`
	LastPoll     int                   `xml:"LAST_POLL"`
	State        State                 `xml:"STATE"`
	LCMState     LCMState              `xml:"LCM_STATE"`
	Resched      int                   `xml:"RESCHED"`
	DeployID     string                `xml:"DEPLOY_ID"`
	Template     Template              `xml:"TEMPLATE"`
	UserTemplate template.UserTemplate `xml:"USER_TEMPLATE"`
	Node         string                `xml:"HISTORY_RECORDS>HISTORY>HOSTNAME"`
}

// Template represents a VM template, i.e. the template structure within VM objects
// (NB: this is not the same as templates in a VM pool).
type Template struct {
	ID     int          `xml:"TEMPLATE_ID"`
	Memory int          `xml:"MEMORY"`
	VMID   int          `xml:"VMID"`
	Disk   []image.Disk `xml:"DISK"`
	CPU    string       `xml:"CPU"`
}

// FromReader reads into a VM struct.
func FromReader(r io.Reader) (*VM, error) {
	v := VM{}
	dec := xml.NewDecoder(r)
	if err := dec.Decode(&v); err != nil {
		return nil, err
	}
	return &v, nil
}

func (v VM) String() string {
	return v.Name
}
