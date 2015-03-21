package h_sillynet

import "sync"

// Threadsafe queue.
type MemoryBlockQueue struct {
	lengthLimit int
	length      int
	items       [][]byte
	locker      sync.Mutex
}

const DefaultMemoryBlockQueueLengthLimit = 1000

func (this *MemoryBlockQueue) Init() {
	if nil == this.items {
		if 0 == this.lengthLimit {
			this.lengthLimit = DefaultMemoryBlockQueueLengthLimit
		}
		this.items = make([][]byte, this.lengthLimit)
	}
}

func (this *MemoryBlockQueue) LengthLimit() int {
	return this.lengthLimit
}

// Returns true if the memory block was successfully pushed & sotred.
// Returns false if there is not enough free space.
func (this *MemoryBlockQueue) Push(memoryBlock []byte) bool {
	this.Init()
	this.locker.Lock()
	var hasFreeSpace = this.length < this.LengthLimit()
	if hasFreeSpace {
		this.items[this.length] = memoryBlock
		this.length = this.length + 1
	}
	this.locker.Unlock()
	return hasFreeSpace
}

func (this *MemoryBlockQueue) Pop() []byte {
	this.Init()
	var shrink1 = func() {
		for i := 0; i < this.length-1; i++ {
			this.items[i] = this.items[i+1]
		}
		this.length = this.length - 1
	}
	var result []byte = nil
	this.locker.Lock()
	var hasItem = this.length > 0
	if hasItem {
		result = this.items[0]
		shrink1()
	}
	this.locker.Unlock()
	return result
}
