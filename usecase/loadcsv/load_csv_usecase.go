package loadcsv

import (
	"CNCMonitorAutoReport/domain/cncmonitor"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

type Executor interface {
	Execute(pic PickingCNCMonitorPram) (CNCMonitorDTO, error)
}

type loadCSVUseCase struct {
	sugar  *zap.SugaredLogger
	loader cncmonitor.Loader
	opener MonitorOpener
}

// Execute implements Executor
func (l loadCSVUseCase) Execute(pic PickingCNCMonitorPram) (CNCMonitorDTO, error) {
	f, err := l.opener.Open(pic.FilePath())
	if err != nil {
		m := fmt.Sprintf("モニタログが開けませんでした path: %s, error: %v", pic.FilePath(), err)
		l.sugar.Error(m)
		return nil, err
	}
	defer f.Close()

	// 文字コード変換
	trans := transform.NewReader(f, japanese.ShiftJIS.NewDecoder())

	cnv, err := pic.convert()
	if err != nil {
		m := fmt.Sprintf("CNC稼働ログ情報変換に失敗しました error: %v", err)
		l.sugar.Error(m)
		return nil, fmt.Errorf(m)
	}

	mon, err := l.loader.Load(trans, cnv)
	if err != nil {
		m := fmt.Sprintf("CNC稼働ログ情報読込に失敗しました error: %v", err)
		l.sugar.Error(m)
		return nil, fmt.Errorf(m)
	}

	mDTO := ParseCNCMonitorDTO(mon)
	return mDTO, nil
}

func New(sugar *zap.SugaredLogger, loader cncmonitor.Loader, opener MonitorOpener) Executor {
	return &loadCSVUseCase{
		sugar:  sugar,
		loader: loader,
		opener: opener,
	}
}

type PickingCNCMonitorPram interface {
	PickingDate() time.Time
	Factory() string
	IPAddress() string
	MachineName() string
	FilePath() string
	convert() (cncmonitor.PickingCNCMonitor, error)
}

type pickingCNCMonitorPram struct {
	pickingDate time.Time
	factory     string
	ipAddress   string
	machineName string
	filePath    string
}

// convert implements PickingCNCMonitorPram
func (p pickingCNCMonitorPram) convert() (cncmonitor.PickingCNCMonitor, error) {
	c, err := cncmonitor.NewPickingCNCMonitor(
		p.PickingDate(),
		p.Factory(),
		p.IPAddress(),
		p.MachineName(),
	)
	if err != nil {
		m := fmt.Sprintf("CNC稼働ログ情報の変換に失敗しました error: %v", err)
		return nil, fmt.Errorf(m)
	}
	return c, nil
}

// FilePath implements PickingCNCMonitorPram
func (p pickingCNCMonitorPram) FilePath() string {
	return p.filePath
}

// Factory implements PickingCNCMonitorPram
func (p pickingCNCMonitorPram) Factory() string {
	return p.factory
}

// IPAddress implements PickingCNCMonitorPram
func (p pickingCNCMonitorPram) IPAddress() string {
	return p.ipAddress
}

// MachineName implements PickingCNCMonitorPram
func (p pickingCNCMonitorPram) MachineName() string {
	return p.machineName
}

// PickingDate implements PickingCNCMonitorPram
func (p pickingCNCMonitorPram) PickingDate() time.Time {
	return p.pickingDate
}

func NewPickingCNCMonitorPram(
	pickingDate time.Time,
	factory string,
	ipAddress string,
	machineName string,
) (PickingCNCMonitorPram, error) {
	f := filepath.Join(
		os.Getenv("BASE_DIR"),
		factory,
		fmt.Sprintf(
			"[%s]_%s.csv",
			ipAddress,
			pickingDate.Format("20060102"),
		))

	return &pickingCNCMonitorPram{
		pickingDate: pickingDate,
		factory:     factory,
		ipAddress:   ipAddress,
		machineName: machineName,
		filePath:    f,
	}, nil
}

// インフラ層・テスト専用
func ReConstructor(
	pickingDate time.Time,
	factory string,
	ipAddress string,
	machineName string,
) PickingCNCMonitorPram {
	obj, _ := NewPickingCNCMonitorPram(
		pickingDate,
		factory,
		ipAddress,
		machineName,
	)
	return obj
}

type CNCMonitorDTO interface {
	Factory() string
	IPAddress() string
	MachineName() string
	Records() []CNCMonitorRecordDTO
}

type cncMonitorDTO struct {
	factory     string
	ipAddress   string
	machineName string
	records     []CNCMonitorRecordDTO
}

// Records implements CNCMonitorDTO
func (m cncMonitorDTO) Records() []CNCMonitorRecordDTO {
	return m.records
}

// Factory implements CNCMonitorDTO
func (m cncMonitorDTO) Factory() string {
	return m.factory
}

// IPAddress implements CNCMonitorDTO
func (m cncMonitorDTO) IPAddress() string {
	return m.ipAddress
}

// MachineName implements CNCMonitorDTO
func (m cncMonitorDTO) MachineName() string {
	return m.machineName
}

func ParseCNCMonitorDTO(mon cncmonitor.CNCMonitorByMachine) CNCMonitorDTO {
	var rs []CNCMonitorRecordDTO
	for _, v := range mon.Records() {
		r := parseCNCMonitorRecordDTO(v)
		rs = append(rs, r)
	}

	return NewCNCMonitorDTO(
		mon.Factory(),
		mon.IPAddress(),
		mon.MachineName(),
		rs,
	)
}

func NewCNCMonitorDTO(
	factory string,
	ipAddress string,
	machineName string,
	records []CNCMonitorRecordDTO,
) CNCMonitorDTO {
	return &cncMonitorDTO{
		factory:     factory,
		ipAddress:   ipAddress,
		machineName: machineName,
		records:     records,
	}
}

type CNCMonitorRecordDTO interface {
	RecordTime() time.Time
	State() string
	ProgramName() string
	FeedRate() int
	SpindleRotation() int
	RunMode() int
	RunState() int
	Emergency() int
}

type cncMonitorRecordDTO struct {
	recordTime      time.Time
	state           string
	programName     string
	feedRate        int
	spindleRotation int
	runMode         int
	runState        int
	emergency       int
}

func parseCNCMonitorRecordDTO(mr cncmonitor.CNCMonitorRecord) CNCMonitorRecordDTO {
	return NewCNCMonitorRecordDTO(
		mr.RecordTime(),
		mr.State(),
		mr.ProgramName(),
		mr.FeedRate(),
		mr.SpindleRotation(),
		mr.RunMode(),
		mr.RunState(),
		mr.Emergency(),
	)
}

// Emergency implements CNCMonitorRecordDTO
func (c cncMonitorRecordDTO) Emergency() int {
	return c.emergency
}

// FeedRate implements CNCMonitorRecordDTO
func (c cncMonitorRecordDTO) FeedRate() int {
	return c.feedRate
}

// ProgramName implements CNCMonitorRecordDTO
func (c cncMonitorRecordDTO) ProgramName() string {
	return c.programName
}

// RecordTime implements CNCMonitorRecordDTO
func (c cncMonitorRecordDTO) RecordTime() time.Time {
	return c.recordTime
}

// RunMode implements CNCMonitorRecordDTO
func (c cncMonitorRecordDTO) RunMode() int {
	return c.runMode
}

// RunState implements CNCMonitorRecordDTO
func (c cncMonitorRecordDTO) RunState() int {
	return c.runState
}

// SpindleRotation implements CNCMonitorRecordDTO
func (c cncMonitorRecordDTO) SpindleRotation() int {
	return c.spindleRotation
}

// State implements CNCMonitorRecordDTO
func (c cncMonitorRecordDTO) State() string {
	return c.state
}

func NewCNCMonitorRecordDTO(
	recordTime time.Time,
	state string,
	programName string,
	feedRate int,
	spindleRotation int,
	runMode int,
	runState int,
	emergency int,
) CNCMonitorRecordDTO {
	return &cncMonitorRecordDTO{
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
