package h_sillynet

import "net"
import "time"

type Client struct {
	messageReceiver  MessageReceiver
	connectionActive bool
	connection       net.Conn
	readerThread     *Thread
	writerThread     *Thread
	incoming         MemoryBlockQueue
	outgoing         MemoryBlockQueue
}

func (this *Client) Push(memoryBlock []byte) {
	this.outgoing.Push(memoryBlock)
}

func (this *Client) Pop() []byte {
	return this.incoming.Pop()
}

func (this *Client) start() {
	var readerThreadRoutine = func(thread *Thread) {
		var readForward = func() {
			this.connection.SetDeadline(time.Now().Add(1 * time.Second))
			var buffer = make([]byte, 1, 1)
			var readLength, readResult = this.connection.Read(buffer)
			if readResult == nil {
				if readLength == 1 {
				}
			}
		}
		for thread.Active {
			if this.connectionActive {
				readForward()
			}
		}
	}
	this.readerThread = StartThread(readerThreadRoutine)
	var writerThreadRoutine = func(thread *Thread) {
	}
	this.writerThread = StartThread(writerThreadRoutine)
}

func (this *Client) Accept(connection net.Conn) {
	this.connection = connection
	this.connectionActive = true
	this.start()
}
