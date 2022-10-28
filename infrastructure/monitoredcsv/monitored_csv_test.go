package monitoredcsv_test

import (
	"CNCMonitorAutoReport/domain/cncmonitor"
	"CNCMonitorAutoReport/infrastructure/monitoredcsv"
	"bufio"
	"io"
	"os"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"

	"go.uber.org/zap"
)

func TestMonitoredCSV_Read(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	sugar := logger.Sugar()

	os.Setenv("BASE_DIR", ".")

	type args struct {
		mon cncmonitor.PickingCNCMonitor
		csv io.Reader
	}
	tests := []struct {
		name    string
		m       cncmonitor.Loader
		args    args
		want    cncmonitor.CNCMonitorByMachine
		wantErr bool
	}{
		{
			name: "正常系_モニタログが読み取れること",
			m:    monitoredcsv.New(sugar),
			args: args{
				mon: csv20221025Monitor(),
				csv: strings.NewReader(csvTestData),
			},
			want:    csv20221025(),
			wantErr: false,
		},
		{
			name: "異常系_モニタログの1行目のフィールド数が2未満のときエラーになること",
			m:    monitoredcsv.New(sugar),
			args: args{
				mon: csv20221025Monitor(),
				csv: strings.NewReader("2未満の行\n2行目\n3行目\n4行目"),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "異常系_モニタログの1行目が日付じゃないときエラーになること",
			m:    monitoredcsv.New(sugar),
			args: args{
				mon: csv20221025Monitor(),
				csv: strings.NewReader("日付,日付じゃない行\n2行目\n3行目\n4行目"),
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.m.Load(tt.args.csv, tt.args.mon)
			if (err != nil) != tt.wantErr {
				t.Errorf("MonitoredCSV.Load() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MonitoredCSV.Load() = %v, want %v", got, tt.want)
			}
		})
	}
}

var pm cncmonitor.PickingCNCMonitor

func csv20221025Monitor() cncmonitor.PickingCNCMonitor {
	if pm == nil {
		pm, _ = cncmonitor.NewPickingCNCMonitor(
			time.Date(2022, 10, 25, 0, 0, 0, 0, time.Local),
			"monitoredcsv",
			"192.168.1.1",
			"CNC-1号機",
		)
	}
	return pm
}

func csv20221025() cncmonitor.CNCMonitorByMachine {
	scn := bufio.NewScanner(strings.NewReader(csvTestData))
	scn.Scan()
	scn.Scan()
	scn.Scan()

	layout := "2006/01/02 15:04"
	var r []cncmonitor.CNCMonitorRecord
	for scn.Scan() {
		v := strings.Split(scn.Text(), ",")
		t, _ := time.ParseInLocation(layout, v[0], time.Local)
		var programName string
		var feedRate, spindleRotation, runMode, runState, emergency int
		if len(v) >= 12 {
			programName = v[2]
			feedRate, _ = strconv.Atoi(v[5])
			spindleRotation, _ = strconv.Atoi(v[6])
			runMode, _ = strconv.Atoi(v[7])
			runState, _ = strconv.Atoi(v[8])
			emergency, _ = strconv.Atoi(v[11])
		}
		r = append(r,
			cncmonitor.NewCNCMonitorRecord(
				t,
				v[1],
				programName,
				feedRate,
				spindleRotation,
				runMode,
				runState,
				emergency,
			))
	}
	mon := cncmonitor.NewCNCMonitorByMachine(
		csv20221025Monitor(),
		r,
	)
	return mon
}
