package cncmonitor_test

import (
	"CNCMonitorAutoReport/domain/cncmonitor"
	"os"
	"testing"
	"time"
)

func TestNewPickingCNCMonitor(t *testing.T) {
	os.Setenv("BASE_DIR", `\\DEVLS\CNCLog\取得データ`)
	type args struct {
		pickingDate time.Time
		factory     string
		ipAddress   string
		machineName string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "正常系_Entity作成後ファイルパスが取得できること",
			args: args{
				pickingDate: time.Date(2022, 10, 20, 0, 0, 0, 0, time.Local),
				factory:     "本社工場",
				ipAddress:   "192.168.11.151",
				machineName: "RB-2号機(本社)",
			},
			want:    `\\DEVLS\CNCLog\取得データ\本社工場\[192.168.11.151]_20221020.csv`,
			wantErr: false,
		},
		{
			name: "異常系_存在しない工場名を渡すとエラーになること",
			args: args{
				pickingDate: time.Now(),
				factory:     "架空",
				ipAddress:   "192.168.11.151",
				machineName: "RB-2号機(本社)",
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "異常系_存在しないIPAddressを渡すとエラーになること",
			args: args{
				pickingDate: time.Now(),
				factory:     "本社工場",
				ipAddress:   "192.168.254.151",
				machineName: "RB-2号機(本社)",
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := cncmonitor.NewPickingCNCMonitor(tt.args.pickingDate, tt.args.factory, tt.args.ipAddress, tt.args.machineName)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewPickingCNCMonitor() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got.FilePath() != tt.want {
				t.Errorf("NewPickingCNCMonitor().FilePath() = %v, want %v", got, tt.want)
			}
		})
	}
}
