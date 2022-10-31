package cncmonitor

import "time"

// CNC稼働レコード
type CNCMonitorRecord interface {
	RecordTime() time.Time
	State() string
	ProgramName() string
	FeedRate() int
	SpindleRotation() int
	RunMode() int
	RunState() int
	Emergency() int
}

type cncMonitorRecord struct {
	recordTime      time.Time
	state           string
	programName     string
	feedRate        int
	spindleRotation int
	runMode         int
	runState        int
	emergency       int
}

// Emergency implements CNCMonitorRecord
func (c cncMonitorRecord) Emergency() int {
	return c.emergency
}

// FeedRate implements CNCMonitorRecord
func (c cncMonitorRecord) FeedRate() int {
	return c.feedRate
}

// ProgramName implements CNCMonitorRecord
func (c cncMonitorRecord) ProgramName() string {
	return c.programName
}

// RecordTime implements CNCMonitorRecord
func (c cncMonitorRecord) RecordTime() time.Time {
	return c.recordTime
}

// RunMode implements CNCMonitorRecord
func (c cncMonitorRecord) RunMode() int {
	return c.runMode
}

// RunState implements CNCMonitorRecord
func (c cncMonitorRecord) RunState() int {
	return c.runState
}

// SpindleRotation implements CNCMonitorRecord
func (c cncMonitorRecord) SpindleRotation() int {
	return c.spindleRotation
}

// State implements CNCMonitorRecord
func (c cncMonitorRecord) State() string {
	return c.state
}

func NewCNCMonitorRecord(
	recordTime time.Time,
	state string,
	programName string,
	feedRate int,
	spindleRotation int,
	runMode int,
	runState int,
	emergency int,
) CNCMonitorRecord {
	return cncMonitorRecord{
		recordTime:      recordTime,
		state:           state,
		programName:     programName,
		feedRate:        feedRate,
		spindleRotation: spindleRotation,
		runMode:         runMode,
		runState:        runState,
		emergency:       emergency,
	}
}
