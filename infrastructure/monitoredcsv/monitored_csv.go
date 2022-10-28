package monitoredcsv

import (
	"CNCMonitorAutoReport/domain/cncmonitor"
	"bufio"
	"fmt"
	"io"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap"
)

type MonitoredCSV struct {
	sugar *zap.SugaredLogger
}

// Read implements cncmonitor.Reader
func (m MonitoredCSV) Load(r io.Reader, mon cncmonitor.PickingCNCMonitor) (cncmonitor.CNCMonitorByMachine, error) {
	file := filepath.Base(mon.FilePath())
	// 1行目の日付取得
	s := bufio.NewScanner(r)
	s.Scan()
	if err := s.Err(); err != nil {
		msg := fmt.Sprintf("モニタログの1行目(日付)が記録されていません file: %s, err: %v", file, err)
		m.sugar.Error(msg)
		return nil, fmt.Errorf(msg)
	}
	_1stLine := strings.Split(s.Text(), ",")
	layout := "2006年01月02日"
	if len(_1stLine) < 2 {
		msg := fmt.Sprintf("モニタログ1行目のフィールド数が不正です line: %v", _1stLine)
		m.sugar.Error(msg)
		return nil, fmt.Errorf(msg)
	}
	recDate, err := time.ParseInLocation(layout, _1stLine[1], time.Local)
	if err != nil {
		msg := fmt.Sprintf("モニタログ1行目の日付形式が不正です line: %s, err: %v", _1stLine[1], err)
		m.sugar.Error(msg)
		return nil, fmt.Errorf(msg)
	}
	// 目標モニタログ日付と記録日付を比較
	if !recDate.Equal(mon.PickingDate()) {
		msg := fmt.Sprintf("モニタログ1行目の日付形式と記録されている日付が違います target: %v, rec: %v", mon.PickingDate(), recDate)
		m.sugar.Error(msg)
		return nil, fmt.Errorf(msg)
	}

	// 2行目の機械名取得
	if !s.Scan() || s.Err() != nil {
		msg := fmt.Sprintf("モニタログに機械名が記録されていません file: %s, err: %v", file, err)
		m.sugar.Error(msg)
		return nil, fmt.Errorf(msg)
	}
	_2ndLine := strings.Split(s.Text(), ",")
	machineName := _2ndLine[1]
	ipAddress := _2ndLine[2]
	m.sugar.Infoln(recDate, machineName, ipAddress)

	var recs []cncmonitor.CNCMonitorRecord
	var i int
	detailLayout := "2006/01/02 15:04"
	for s.Scan() {
		i++
		if i == 1 {
			continue
		}
		// 列数が可変長のためcsv,gota/dataframeパッケージが使いにくい Splitを使用する
		_line := strings.Split(s.Text(), ",")
		d, err := time.ParseInLocation(detailLayout, _line[0], time.Local)
		if err != nil {
			msg := fmt.Sprintf("モニタログ%d行目の日付形式が不正です line: %s, err: %v", i, _line[0], err)
			m.sugar.Error(msg)
			return nil, fmt.Errorf(msg)
		}

		state := _line[1]

		var (
			programName                                             string
			feedRate, spindleRotation, runMode, runState, emergency int
		)

		if len(_line) >= 12 {
			programName = _line[2]
			feedRate, _ = strconv.Atoi(_line[5])
			spindleRotation, _ = strconv.Atoi(_line[6])
			runMode, _ = strconv.Atoi(_line[7])
			runState, _ = strconv.Atoi(_line[8])
			emergency, _ = strconv.Atoi(_line[11])
		}

		rec := cncmonitor.NewCNCMonitorRecord(
			d,
			state,
			programName,
			feedRate,
			spindleRotation,
			runMode,
			runState,
			emergency,
		)
		recs = append(recs, rec)
	}

	return cncmonitor.NewCNCMonitorByMachine(
		mon,
		recs,
	), nil
}

func New(sugar *zap.SugaredLogger) cncmonitor.Loader {
	return &MonitoredCSV{
		sugar: sugar,
	}
}
