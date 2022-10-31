//go:generate mockgen -source=$GOFILE -destination=mock_$GOPACKAGE/mock_$GOFILE
package cncmonitor

import "io"

type Loader interface {
	Load(r io.Reader, mon PickingCNCMonitor) (CNCMonitorByMachine, error)
}
