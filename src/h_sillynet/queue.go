package h_sillynet

import "sync"

// Threadsafe queue.
type MemoryBlockQueue struct {
	length int
	items  [][]byte
	locker sync.Locker
}

func NewMemoryBlockQueue(lengthLimit int) *MemoryBlockQueue {
	return &MemoryBlockQueue{items: make([][]byte, lengthLimit)}
}

func (this *MemoryBlockQueue) LengthLimit() int {
	return len(this.items)
}

func (this *MemoryBlockQueue) Push(memoryBlock []byte) bool {
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
