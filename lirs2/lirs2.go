package lirs2

import (
	"github.com/secnot/orderedmap"
	"lirs2/simulator"
	"os"
	"time"
)

type Instance struct {
	orderedInstance1 *orderedmap.OrderedMap
	orderedInstance2 *orderedmap.OrderedMap
}
type LIRS2 struct {
	cacheSize    int
	hit          int
	miss         int
	writeCount   int
	LIRSize      int
	HIRSize      int
	orderedStack Instance
	orderedList  *orderedmap.OrderedMap
	LIR          map[interface{}]int
	HIR          map[interface{}]int
	cache        map[interface{}]bool
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
