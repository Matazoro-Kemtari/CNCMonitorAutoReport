package loadcsv_test

import (
	"CNCMonitorAutoReport/domain/cncmonitor/mock_cncmonitor"
	"CNCMonitorAutoReport/domain/cncmonitor/testfactory"
	"CNCMonitorAutoReport/usecase/loadcsv"
	"CNCMonitorAutoReport/usecase/loadcsv/mock_loadcsv"
	"os"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"go.uber.org/zap"
)

func Test_loadCSVUseCase_Execute(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	sugar := logger.Sugar()

	// モックコントローラーの生成
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// cncmonitor.Loaderモック
	mock_ldr := mock_cncmonitor.NewMockLoader(ctrl)
	res := testfactory.TestingMonitorCreate()
	mock_ldr.EXPECT().Load(gomock.Any(), gomock.Any()).Return(res, nil)

	// MonitorOpenerモック
	mock_opener := mock_loadcsv.NewMockMonitorOpener(ctrl)
	f, _ := os.Create("test")
	defer func() {
		f.Close()
		os.Remove("test")
	}()

	mock_opener.EXPECT().Open(gomock.Any()).Return(f, nil)

	type args struct {
		pic loadcsv.PickingCNCMonitorPram
	}
	tests := []struct {
		name    string
		l       loadcsv.Executor
		args    args
		want    loadcsv.CNCMonitorDTO
		wantErr bool
	}{
		{
			name: "正常系_モニタログが読み込めること",
			l:    loadcsv.New(sugar, mock_ldr, mock_opener),
			args: args{
				pic: loadcsv.ReConstructor(
					res.PickingDate(),
					res.Factory(),
					res.IPAddress(),
					res.MachineName(),
				),
			},
			want:    loadcsv.ParseCNCMonitorDTO(res),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.l.Execute(tt.args.pic)
			if (err != nil) != tt.wantErr {
				t.Errorf("loadCSVUseCase.Execute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("loadCSVUseCase.Execute() = %v, want %v", got, tt.want)
			}
		})
	}
}
