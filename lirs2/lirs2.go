package lirs2

import (
	"github.com/secnot/orderedmap"
	"lirs2/simulator"
	"os"
	"time"
)

type Instance struct {
	blockNumber int
	isHot       bool
	isInstance1 bool
	isResident  bool
}

type LIRS2 struct {
	cacheSize  int
	hit        int
	miss       int
	writeCount int
	LIRSize    int
	HIRSize    int
	LIRS2Queue *orderedmap.OrderedMap
	CoReQueue  *orderedmap.OrderedMap
	LIR        map[interface{}]int
	HIR        map[interface{}]int
}

func NewLIRS2(cacheSize int, HIRSize int) *LIRS2 {
	if HIRSize > 100 || HIRSize < 0 {
		panic("HIRSize must be between 0 and 100")
	}
	LIRCapacity := (100 - HIRSize) * cacheSize / 100
	HIRCapacity := HIRSize * cacheSize / 100
	return &LIRS2{
		cacheSize:  cacheSize,
		hit:        0,
		miss:       0,
		writeCount: 0,
		LIRSize:    LIRCapacity,
		HIRSize:    HIRCapacity,
		LIRS2Queue: orderedmap.NewOrderedMap(),
		CoReQueue:  orderedmap.NewOrderedMap(),
		LIR:        make(map[interface{}]int, LIRCapacity),
		HIR:        make(map[interface{}]int, HIRCapacity),
	}
}

func (LIRS2Object *LIRS2) PrintToFile(file *os.File, start time.Time) error {
	//TODO implement me
	panic("implement me")
}

func (LIRS2Object *LIRS2) Get(trace simulator.Trace) error {
	//TODO implement me
	instance := Instance{
		blockNumber: trace.Address,
		isHot:       false,
		isInstance1: false,
		isResident:  false,
	}

	blockAddress := trace.Address
	operation := trace.Operation
	if operation == "W" {
		LIRS2Object.writeCount++
	}
	if len(LIRS2Object.LIR) < LIRS2Object.LIRSize {
		// LIR is not full; there is space in cache
		LIRS2Object.miss += 1
		if _, ok := LIRS2Object.LIR[blockAddress]; ok {
			// block is in LIR, not a miss
			LIRS2Object.miss -= 1
			LIRS2Object.hit += 1
		}
		LIRS2Object.addToStack(instance)
		LIRS2Object.makeLIR(blockAddress)
		return nil
	}
	if _, ok := LIRS2Object.LIR[blockAddress]; ok {
		// block is in LIR, hit
		LIRS2Object.handleLIRBlock(blockAddress)
	} else if _, ok := LIRS2Object.HIR[blockAddress]; ok {
		// block is in HIR, hit
		LIRS2Object.handleHIRResidentBlock(blockAddress)
	} else {
		// block is not in LIR or HIR Resident, miss
		LIRS2Object.handleHIRNonResidentBlock(blockAddress)
	}
	return nil
}

func (LIRS2Object *LIRS2) addToStack(instance Instance) {
	//TODO implement me
	LIRS2Object.LIRS2Queue.Set(instance, 1)
	key, _, ok := LIRS2Object.LIRS2Queue.GetFirst()
	if !ok {
		panic("LIRS2Queue is empty")
	}
	if key == instance {
		//TODO implement me
	}
	panic("implement me")
}

func (LIRS2Object *LIRS2) makeLIR(block int) {
	//TODO implement me
	panic("implement me")
}

func (LIRS2Object *LIRS2) handleLIRBlock(block int) {
	LIRS2Object.hit += 1
}

func (LIRS2Object *LIRS2) handleHIRResidentBlock(block int) {
	//TODO implement me
	panic("implement me")
}

func (LIRS2Object *LIRS2) handleHIRNonResidentBlock(block int) {
	//TODO implement me
	panic("implement me")
}
