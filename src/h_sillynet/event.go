package h_sillynet

import "sync"

// Wrapper for sync.Cond type.
func NewEvent() Event {
	var mutex = &sync.Mutex{}
	var result = Event{}
	result.condition = sync.NewCond(mutex)
	return result
}

var EmptyEvent = Event{}

// Create with constructor only (NewEvent).
type Event struct {
	condition *sync.Cond
}

func (this Event) Wait() {
	this.condition.L.Lock()
	this.condition.Wait()
	this.condition.L.Unlock()
}

func (this Event) Signal() {
	this.condition.L.Lock()
	this.condition.Signal()
	this.condition.L.Unlock()
}

func (this Event) Exists() bool {
	return this.condition != nil
}
