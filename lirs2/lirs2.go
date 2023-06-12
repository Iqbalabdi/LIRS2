package lirs2

import (
	"errors"
	"github.com/secnot/orderedmap"
	"lirs2/simulator"
	"os"
	"time"
)

type (
	Instance struct {
		block      int
		accessTime int64
	}

	LIRS2 struct {
		cacheSize      int
		hit            int
		miss           int
		writeCount     int
		LIRSize        int
		HIRSize        int
		Instance1Queue *orderedmap.OrderedMap
		Instance2Queue *orderedmap.OrderedMap
		CoReQueue      *orderedmap.OrderedMap
		LIRBlock       map[interface{}]int
		HIRBlock       map[interface{}]int
	}
)

func NewLIRS2(cacheSize int, HIRSize int) *LIRS2 {
	if HIRSize > 100 || HIRSize < 0 {
		panic("HIRSize must be between 0 and 100")
	}
	LIRCapacity := (100 - HIRSize) * cacheSize / 100
	HIRCapacity := HIRSize * cacheSize / 100
	return &LIRS2{
		cacheSize:      cacheSize,
		hit:            0,
		miss:           0,
		writeCount:     0,
		LIRSize:        LIRCapacity,
		HIRSize:        HIRCapacity,
		Instance1Queue: orderedmap.NewOrderedMap(),
		Instance2Queue: orderedmap.NewOrderedMap(),
		CoReQueue:      orderedmap.NewOrderedMap(),
		LIRBlock:       make(map[interface{}]int, LIRCapacity),
		HIRBlock:       make(map[interface{}]int, HIRCapacity),
	}
}

func (LIRS2Object *LIRS2) PrintToFile(file *os.File, start time.Time) error {
	//TODO implement me
	panic("implement me")
}

func (LIRS2Object *LIRS2) Get(trace simulator.Trace) error {
	//init data
	data := &Instance{
		block:      trace.Address,
		accessTime: time.Now().Unix(),
	}

	operation := trace.Operation
	if operation == "W" {
		LIRS2Object.writeCount++
	}
	if len(LIRS2Object.LIRBlock) < LIRS2Object.LIRSize {
		// LIRBlock is not full; there is space in cache
		LIRS2Object.miss += 1
		if _, ok := LIRS2Object.LIRBlock[data.block]; ok {
			// block is in LIRBlock, not a miss
			LIRS2Object.miss -= 1
			LIRS2Object.hit += 1
		}
		LIRS2Object.makeLIR(data)
		return nil
	}

	if _, ok := LIRS2Object.LIRBlock[data.block]; ok {
		// block is in LIRBlock, hit
		LIRS2Object.handleLIRBlock(data)
	} else if _, ok := LIRS2Object.CoReQueue.Get(data.block); ok {
		// block is in HIRBlock, hit
		LIRS2Object.handleHIRResidentBlock(data)
	} else {
		// block is not in LIRBlock or HIRBlock Resident, miss
		LIRS2Object.handleHIRNonResidentBlock(data)
	}
	return nil
}

func (LIRS2Object *LIRS2) makeLIR(data *Instance) {
	LIRS2Object.LIRBlock[data.block] = 1
	LIRS2Object.removeFromCoreQueue(data.block)
	delete(LIRS2Object.HIRBlock, data.block)
}

func (LIRS2Object *LIRS2) handleLIRBlock(data *Instance) {
	LIRS2Object.hit += 1
	if _, ok := LIRS2Object.Instance2Queue.Get(data.block); ok {
		if key, _, _ := LIRS2Object.Instance2Queue.GetFirst(); key == data.block {
			LIRS2Object.stackPruning(false)
		}
		LIRS2Object.Instance2Queue.Delete(data.block)
	}
	if _, ok := LIRS2Object.Instance1Queue.Get(data.block); ok {
		LIRS2Object.changeToInstance2(data)
	}
	LIRS2Object.Instance1Queue.Set(data.block, data)
}

func (LIRS2Object *LIRS2) handleHIRResidentBlock(data *Instance) {
	LIRS2Object.hit += 1
	if _, ok := LIRS2Object.Instance2Queue.Get(data.block); ok {
		LIRS2Object.makeLIR(data)
		LIRS2Object.removeFromCoreQueue(data.block)
		LIRS2Object.stackPruning(true)
	}
	if _, ok := LIRS2Object.Instance1Queue.Get(data.block); ok {
		LIRS2Object.changeToInstance2(data)
	}
	LIRS2Object.CoReQueue.MoveLast(data.block)
	LIRS2Object.Instance1Queue.Set(data.block, data)
}

func (LIRS2Object *LIRS2) handleHIRNonResidentBlock(data *Instance) {
	LIRS2Object.miss += 1
	if _, ok := LIRS2Object.Instance2Queue.Get(data.block); ok {
		LIRS2Object.makeLIR(data)
		LIRS2Object.stackPruning(true)
	}
	if _, ok := LIRS2Object.Instance1Queue.Get(data.block); ok {
		LIRS2Object.changeToInstance2(data)
	}
	LIRS2Object.makeHIR(data)
	LIRS2Object.addToCoreQueue(data.block)
	LIRS2Object.Instance1Queue.Set(data.block, data)
}

func (LIRS2Object *LIRS2) removeFromCoreQueue(block int) {
	LIRS2Object.CoReQueue.Delete(block)
}

func (LIRS2Object *LIRS2) makeHIR(data *Instance) {
	LIRS2Object.HIRBlock[data.block] = 1
	delete(LIRS2Object.LIRBlock, data.block)
}

func (LIRS2Object *LIRS2) changeToInstance2(data *Instance) {
	LIRS2Object.Instance1Queue.Delete(data.block)
	LIRS2Object.Instance2Queue.Set(data.block, data)
}

func (LIRS2Object *LIRS2) addToCoreQueue(block int) {
	if LIRS2Object.CoReQueue.Len() < LIRS2Object.HIRSize {
		LIRS2Object.CoReQueue.PopFirst()
	}
	LIRS2Object.CoReQueue.Set(block, 1)
}

func (LIRS2Object *LIRS2) stackPruning(removeLIR bool) error {
	_, val, ok := LIRS2Object.Instance2Queue.GetFirst()
	if !ok {
		return errors.New("Instance2Queue is empty")
	}

	if removeLIR {
		LIRS2Object.makeHIR(val.(*Instance))
		LIRS2Object.addToCoreQueue(val.(*Instance).block)
	}

	iter := LIRS2Object.Instance2Queue.Iter()
	for k, _, ok := iter.Next(); ok; k, _, ok = iter.Next() {
		if LIRS2Object.LIRBlock[k] == 1 {
			break
		}
		LIRS2Object.Instance2Queue.Delete(key)
	}
	iter = LIRS2Object.Instance1Queue.Iter()
	for _, v, ok := iter.Next(); ok; _, v, ok = iter.Next() {
		if val.(*Instance).accessTime < v.(*Instance).accessTime {

		}
	}
	return nil
}
