package lirs2

import (
	"errors"
	"fmt"
	"lirs2/pkg/orderedmap"
	"lirs2/simulator"
	"os"
	"time"
)

type (
	Instance struct {
		block       int
		accessCount int
	}

	LIRS2 struct {
		accessCounter   int
		cacheSize       int
		hit             int
		miss            int
		writeCount      int
		readCount       int
		LIRSize         int
		HIRSize         int
		Instance1Queue  *orderedmap.OrderedMap
		Instance2Queue  *orderedmap.OrderedMap
		CoReQueue       *orderedmap.OrderedMap
		LIRBlock        map[interface{}]int
		HIRBlock        map[interface{}]int
		Instance2Access map[interface{}]int
	}
)

/*var InfoLogger *log.Logger

func init() {
	f, err := os.OpenFile("logfile", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}

	InfoLogger = log.New(f, "INFO: ", log.LstdFlags)
}*/

func NewLIRS2(cacheSize int, HIRSize int) *LIRS2 {
	if HIRSize > 100 || HIRSize < 0 {
		panic("HIRSize must be between 0 and 100")
	}
	LIRCapacity := (100 - HIRSize) * cacheSize / 100
	HIRCapacity := HIRSize * cacheSize / 100
	return &LIRS2{
		accessCounter:   0,
		cacheSize:       cacheSize,
		hit:             0,
		miss:            0,
		writeCount:      0,
		readCount:       0,
		LIRSize:         LIRCapacity,
		HIRSize:         HIRCapacity,
		Instance1Queue:  orderedmap.NewOrderedMap(),
		Instance2Queue:  orderedmap.NewOrderedMap(),
		CoReQueue:       orderedmap.NewOrderedMap(),
		LIRBlock:        make(map[interface{}]int, LIRCapacity),
		HIRBlock:        make(map[interface{}]int, HIRCapacity),
		Instance2Access: map[interface{}]int{},
	}
}

func (LIRS2Object *LIRS2) Get(trace simulator.Trace) error {
	//init data
	LIRS2Object.accessCounter++
	data := &Instance{
		block:       trace.Address,
		accessCount: LIRS2Object.accessCounter,
	}

	operation := trace.Operation
	if operation == "W" {
		LIRS2Object.writeCount++
	} else {
		LIRS2Object.readCount++
	}
	if len(LIRS2Object.LIRBlock) < LIRS2Object.LIRSize {
		// LIRBlock is not full; there is space in cache
		LIRS2Object.miss += 1
		if _, ok := LIRS2Object.LIRBlock[data.block]; ok {
			// block is in LIRBlock, not a miss
			LIRS2Object.miss -= 1
			LIRS2Object.hit += 1
			if _, ok := LIRS2Object.Instance2Queue.Get(data.block); ok {
				LIRS2Object.Instance2Queue.Delete(data.block)
			}
			if _, ok := LIRS2Object.Instance1Queue.Get(data.block); ok {
				LIRS2Object.changeToInstance2(data)
			}
		} else {
			LIRS2Object.makeLIR(data)
		}
		LIRS2Object.Instance1Queue.Set(data.block, data)
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

func (LIRS2Object *LIRS2) handleLIRBlock(data *Instance) {
	LIRS2Object.hit += 1
	if _, ok := LIRS2Object.Instance2Queue.Get(data.block); ok {
		if key, _, _ := LIRS2Object.Instance2Queue.GetFirst(); key.(int) == data.block {
			LIRS2Object.stackPruning(false)
		} else {
			LIRS2Object.Instance2Queue.Delete(data.block)
		}
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
		LIRS2Object.Instance2Queue.Delete(data.block)
		LIRS2Object.removeFromCoreQueue(data.block)
		LIRS2Object.stackPruning(true)
	} else {
		LIRS2Object.CoReQueue.MoveLast(data.block)
	}
	if _, ok := LIRS2Object.Instance1Queue.Get(data.block); ok {
		LIRS2Object.changeToInstance2(data)
	}
	LIRS2Object.Instance1Queue.Set(data.block, data)
}

func (LIRS2Object *LIRS2) handleHIRNonResidentBlock(data *Instance) {
	LIRS2Object.miss += 1
	if _, ok := LIRS2Object.Instance2Queue.Get(data.block); ok {
		LIRS2Object.makeLIR(data)
		LIRS2Object.Instance2Queue.Delete(data.block)
		LIRS2Object.stackPruning(true)

	} else {
		LIRS2Object.makeHIR(data)
		LIRS2Object.addToCoreQueue(data.block)
	}
	if _, ok := LIRS2Object.Instance1Queue.Get(data.block); ok {
		LIRS2Object.changeToInstance2(data)
	}
	LIRS2Object.Instance1Queue.Set(data.block, data)
}

func (LIRS2Object *LIRS2) makeLIR(data *Instance) {
	LIRS2Object.LIRBlock[data.block] = 1
	LIRS2Object.removeFromCoreQueue(data.block)
	delete(LIRS2Object.HIRBlock, data.block)
}

func (LIRS2Object *LIRS2) makeHIR(data *Instance) {
	LIRS2Object.HIRBlock[data.block] = 1
	delete(LIRS2Object.LIRBlock, data.block)
}

func (LIRS2Object *LIRS2) changeToInstance2(data *Instance) {
	val, _ := LIRS2Object.Instance1Queue.Get(data.block)
	LIRS2Object.Instance1Queue.Delete(val.(*Instance).block)
	LIRS2Object.Instance2Queue.Set(val.(*Instance).block, val)

	LIRS2Object.Instance2Access[val.(*Instance).accessCount] = val.(*Instance).block
	for i := val.(*Instance).accessCount + 1; i <= LIRS2Object.accessCounter; i++ {
		if block, ok := LIRS2Object.Instance2Access[i]; ok {
			LIRS2Object.Instance2Queue.MoveToSpecificIndex(val.(*Instance).block, block)
			break
		}
	}
	//iter := LIRS2Object.Instance2Queue.Iter()
	//for _, v, ok := iter.Next(); ok; _, v, ok = iter.Next() {
	//	if val.(*Instance).accessCount < v.(*Instance).accessCount {
	//		//fmt.Println("move to specific index and access Count", val.(*Instance).accessCount, v.(*Instance).accessCount)
	//		LIRS2Object.Instance2Queue.MoveToSpecificIndex(val.(*Instance).block, v.(*Instance).block)
	//		//fmt.Println("done")
	//		break
	//	}
	//}
}

func (LIRS2Object *LIRS2) addToCoreQueue(block int) {
	if LIRS2Object.CoReQueue.Len() == LIRS2Object.HIRSize {
		LIRS2Object.CoReQueue.PopFirst()
	}
	LIRS2Object.CoReQueue.Set(block, 1)
}

func (LIRS2Object *LIRS2) removeFromCoreQueue(block int) {
	LIRS2Object.CoReQueue.Delete(block)
}

func (LIRS2Object *LIRS2) stackPruning(removeLIR bool) error {
	var flag *Instance

	_, val, ok := LIRS2Object.Instance2Queue.PopFirst()
	if !ok {
		return errors.New("Instance2Queue is empty")
	}

	if removeLIR {
		LIRS2Object.makeHIR(val.(*Instance))
		LIRS2Object.addToCoreQueue(val.(*Instance).block)
	}

	// delete instance2 in queue if it is not LIR
	iter := LIRS2Object.Instance2Queue.Iter()
	for k, v, ok := iter.Next(); ok; k, v, ok = iter.Next() {
		if _, ok := LIRS2Object.LIRBlock[k]; ok {
			flag = v.(*Instance)
			break
		}
		delete(LIRS2Object.Instance2Access, v.(*Instance).accessCount)
		LIRS2Object.Instance2Queue.PopFirst()
	}

	// delete instance1 in queue if access-time is less than bottom instance2
	iter = LIRS2Object.Instance1Queue.Iter()
	for _, v, ok := iter.Next(); ok; _, v, ok = iter.Next() {
		if flag != nil && flag.accessCount < v.(*Instance).accessCount {
			break
		}
		if _, ok := LIRS2Object.LIRBlock[v.(*Instance).block]; ok {
			delete(LIRS2Object.LIRBlock, v.(*Instance).block)
		}
		LIRS2Object.Instance1Queue.PopFirst()
	}
	return nil
}

func (LIRS2Object *LIRS2) PrintToFile(file *os.File, start time.Time) error {
	timeExec := time.Since(start)
	hitRatio := float32(LIRS2Object.hit) / float32(LIRS2Object.hit+LIRS2Object.miss) * 100
	result := fmt.Sprintf(`-----------------------------------------------------
LIRS2
cache size : %v
cache hit : %v
cache miss : %v
hit ratio : %v
Instance2Queue : %v
Instance1Queue : %v
CoreQueue size : %v
LIR capacity: %v
HIR capacity: %v
write count : %v
read count : %v
time execution : %v
access count : %v
`, LIRS2Object.cacheSize, LIRS2Object.hit, LIRS2Object.miss, hitRatio, LIRS2Object.Instance2Queue.Len(), LIRS2Object.Instance1Queue.Len(),
		LIRS2Object.CoReQueue.Len(), LIRS2Object.LIRSize, LIRS2Object.HIRSize, LIRS2Object.writeCount, LIRS2Object.readCount, timeExec.Seconds(), LIRS2Object.accessCounter)
	_, err := file.WriteString(result)
	return err
}
