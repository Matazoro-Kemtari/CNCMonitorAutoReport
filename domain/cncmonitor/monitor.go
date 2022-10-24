package cncmonitor

// CNC稼働ログ
type CNCMonitorByMachine interface {
	MachineName() string
	Records() []CNCMonitorRecord
}

type cncMonitorByMachine struct {
	machineName string
	records     []CNCMonitorRecord
}

// Records implements CNCMonitor
func (c cncMonitorByMachine) Records() []CNCMonitorRecord {
	return c.records
}

// MachineName implements CNCMonitor
func (c cncMonitorByMachine) MachineName() string {
	return c.machineName
}

func NewCNCMonitorByMachine(machineName string, records []CNCMonitorRecord) CNCMonitorByMachine {
	return &cncMonitorByMachine{
		machineName: machineName,
		records:     records,
	}
}
