package monitoredcsv

import (
	"CNCMonitorAutoReport/usecase/loadcsv"
	"fmt"
	"os"

	"go.uber.org/zap"
)

type MonitoredLog struct {
	sugar *zap.SugaredLogger
}

// Open implements loadcsv.MonitorOpener
func (m MonitoredLog) Open(p string) (*os.File, error) {
	// モニタログを開く
	f, err := os.Open(p)
	if err != nil {
		msg := fmt.Sprintf("CNC稼働ログが開けませんでした path: %s, error: %v", p, err)
		m.sugar.Error(msg)
		return nil, fmt.Errorf(msg)
	}

	return f, nil
}

func NewMonitoredLog(sugar *zap.SugaredLogger) loadcsv.MonitorOpener {
	return &MonitoredLog{
		sugar: sugar,
	}
}
