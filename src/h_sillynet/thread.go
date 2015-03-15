package h_sillynet

import "sync"

type Thread struct {
	waitGroup sync.WaitGroup
	Active    bool
}

func StartThread(f func(thread *Thread)) *Thread {
	var thread = &Thread{}
	thread.waitGroup.Add(1)
	thread.Active = true
	go func() {
		f(thread)
		thread.waitGroup.Done()
	}()
	return thread
}

func (this *Thread) WaitFor() {
	this.waitGroup.Wait()
}
