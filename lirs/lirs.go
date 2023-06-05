package lirs

import (
	"lirs2/simulator"
	"os"
	"time"
)

type LIRS struct {
}

func NewLIRS(cacheSize int) *LIRS {
	panic("Not implemented")
}

func (LIRSObject LIRS) Get(trace simulator.Trace) error {
	//TODO implement me
	panic("implement me")
}

func (LIRSObject LIRS) PrintToFile(file *os.File, start time.Time) error {
	//TODO implement me
	panic("implement me")
}
