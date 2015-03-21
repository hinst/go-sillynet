package h_sillynet

import "net"
import "time"

type Client struct {
	messageReceiver MessageReceiver
	connection      net.Conn
	readerThread    *Thread
	writerThread    *Thread
	incoming        MemoryBlockQueue
	outgoing        MemoryBlockQueue
}

func (this *Client) Push(memoryBlock []byte) {
	this.outgoing.Push(memoryBlock)
}

func (this *Client) Pop() []byte {
	return this.incoming.Pop()
}

func CheckIfNetTimeoutError(e error) bool {
	var netError, castResult = readResult.(net.Error)
	if castResult {
		result = netError.Timeout()
	} else {
		result = false
	}
}

func (this *Client) start() {
	var buffer = make([]byte, 1, 1)
	var tryExtractMessage = func() {
		var incomingMessage = this.messageReceiver.Extract()
		if incomingMessage != nil {
			this.incoming.Push(incomingMessage)
		}
	}
	var readForward = func() {
		this.connection.SetDeadline(time.Now().Add(1 * time.Second))
		var readLength, readResult = this.connection.Read(buffer)
		if readResult == nil {
			if readLength == 1 {
				this.messageReceiver.Write(buffer[0])
				tryExtractMessage()
			}
		} else if false == CheckIfNetTimeoutError(readResult) {
			this.connection.Close()
			this.connection = nil
		}
	}
	var readerThreadRoutine = func(thread *Thread) {
		for thread.Active {
			if this.ConnectionActive() {
				readForward()
			}
		}
	}
	this.readerThread = StartThread(readerThreadRoutine)
	var writerThreadRoutine = func(thread *Thread) {
	}
	this.writerThread = StartThread(writerThreadRoutine)
}

func (this *Client) Accept(connection net.Conn) bool {
	if this.connection == nil {
		this.connection = connection
		this.connectionActive = true
		this.start()
		return true
	} else {
		return false
	}
}

func (this *Client) ConnectionActive() bool {
	return this.connection != nil
}
