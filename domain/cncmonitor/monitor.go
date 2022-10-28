package cncmonitor

// CNC稼働ログ
type CNCMonitorByMachine interface {
	Records() []CNCMonitorRecord
}

type cncMonitorByMachine struct {
	PickingCNCMonitor
	records []CNCMonitorRecord
}

// Records implements CNCMonitor
func (c cncMonitorByMachine) Records() []CNCMonitorRecord {
	return c.records
}

func NewCNCMonitorByMachine(
	pickingCNCMonitor PickingCNCMonitor,
	records []CNCMonitorRecord,
) CNCMonitorByMachine {
	return &cncMonitorByMachine{
		PickingCNCMonitor: pickingCNCMonitor,
		records:           records,
	}
}
