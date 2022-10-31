package testfactory

import (
	"CNCMonitorAutoReport/domain/cncmonitor"
	"time"
)

type pickingCNCMonitorOptions struct {
	pickingDate time.Time
	factory     string
	ipAddress   string
	machineName string
	filePath    string
}

type pickingCNCMonitorOption func(*pickingCNCMonitorOptions)

func TestingPickingCNCMonitorCreate(options ...pickingCNCMonitorOption) cncmonitor.PickingCNCMonitor {
	// デフォルト値
	opts := &pickingCNCMonitorOptions{
		pickingDate: time.Now(),
		factory:     "本社",
		ipAddress:   "192.168.1.1",
		machineName: "CNC-1号機",
		filePath:    `./[192.168.1.1]_20220101.csv`,
	}

	for _, option := range options {
		option(opts)
	}

	obj, _ := cncmonitor.NewPickingCNCMonitor(
		opts.pickingDate,
		opts.factory,
		opts.ipAddress,
		opts.machineName,
	)
	return obj
}

type cncMonitorByMachineOptions struct {
	pickingCNCMonitor cncmonitor.PickingCNCMonitor
	records           []cncmonitor.CNCMonitorRecord
}

type cncMonitorByMachineOption func(*cncMonitorByMachineOptions)

func TestingMonitorCreate(options ...cncMonitorByMachineOption) cncmonitor.CNCMonitorByMachine {
	// デフォルト値
	opts := &cncMonitorByMachineOptions{
		pickingCNCMonitor: TestingPickingCNCMonitorCreate(),
		records:           []cncmonitor.CNCMonitorRecord{TestingRecordCreate()},
	}

	for _, option := range options {
		option(opts)
	}

	return cncmonitor.NewCNCMonitorByMachine(
		opts.pickingCNCMonitor,
		opts.records,
	)
}

type cncMonitorRecordOptions struct {
	recordTime      time.Time
	state           string
	programName     string
	feedRate        int
	spindleRotation int
	runMode         int
	runState        int
	emergency       int
}

type cncMonitorRecordOption func(*cncMonitorRecordOptions)

func TestingRecordCreate(options ...cncMonitorRecordOption) cncmonitor.CNCMonitorRecord {
	// デフォルト値
	opts := &cncMonitorRecordOptions{
		recordTime:      time.Now(),
		state:           "接続（電源ON)",
		programName:     "O0",
		feedRate:        5,
		spindleRotation: 7,
		runMode:         1,
		runState:        3,
		emergency:       0,
	}

	for _, option := range options {
		option(opts)
	}

	return cncmonitor.NewCNCMonitorRecord(
		opts.recordTime,
		opts.state,
		opts.programName,
		opts.feedRate,
		opts.spindleRotation,
		opts.runMode,
		opts.runState,
		opts.emergency,
	)
}
