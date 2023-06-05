package lirs2

import (
	"lirs2/simulator"
	"os"
	"time"
)

type LIRS2 struct {
}

func NewLIRS2(cacheSize int) *LIRS2 {
	panic("Not implemented")
}

func (LIRS2Object LIRS2) Get(trace simulator.Trace) error {
	//TODO implement me
	panic("implement me")
}

func (LIRS2Object LIRS2) PrintToFile(file *os.File, start time.Time) error {
	//TODO implement me
	panic("implement me")
}
