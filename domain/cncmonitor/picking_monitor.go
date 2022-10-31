package cncmonitor

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type PickingCNCMonitor interface {
	PickingDate() time.Time
	Factory() string
	IPAddress() string
	MachineName() string
	FilePath() string
}

type pickingCNCMonitor struct {
	pickingDate time.Time
	factory     string
	ipAddress   string
	machineName string
	filePath    string
}

// FilePath implements PickingCNCMonitor
func (p pickingCNCMonitor) FilePath() string {
	return p.filePath
}

// Factory implements PickingCNCMonitor
func (p pickingCNCMonitor) Factory() string {
	return p.factory
}

// IPAddress implements PickingCNCMonitor
func (p pickingCNCMonitor) IPAddress() string {
	return p.ipAddress
}

// MachineName implements PickingCNCMonitor
func (p pickingCNCMonitor) MachineName() string {
	return p.machineName
}

// PickingDate implements PickingCNCMonitor
func (p pickingCNCMonitor) PickingDate() time.Time {
	return p.pickingDate
}

func NewPickingCNCMonitor(
	pickingDate time.Time,
	factory string,
	ipAddress string,
	machineName string,
) (PickingCNCMonitor, error) {
	f := filepath.Join(
		os.Getenv("BASE_DIR"),
		factory,
		fmt.Sprintf(
			"[%s]_%s.csv",
			ipAddress,
			pickingDate.Format("20060102"),
		))

	return &pickingCNCMonitor{
		pickingDate: pickingDate,
		factory:     factory,
		ipAddress:   ipAddress,
		machineName: machineName,
		filePath:    f,
	}, nil
}
