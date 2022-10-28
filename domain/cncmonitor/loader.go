package cncmonitor

import "io"

type Loader interface {
	Load(r io.Reader, mon PickingCNCMonitor) (CNCMonitorByMachine, error)
}
