//go:generate mockgen -source=$GOFILE -destination=mock_$GOPACKAGE/mock_$GOFILE
package loadcsv

import "os"

type MonitorOpener interface {
	Open(p string) (*os.File, error)
}
