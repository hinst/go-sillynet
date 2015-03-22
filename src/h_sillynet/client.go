package h_sillynet

import "net"
import "time"
import "sync"

func NewClient() *Client {
	var result = &Client{}
	result.ReaderThreadThrashingInterval = DefaultClientReaderThreadThrashingInterval()
	result.WriterThreadThrashingInterval = DefaultClientWriterThreadThrashingInterval()
	return result
}

type Client struct {
	connection       net.Conn
	connectionLocker sync.Mutex

	readerThread *Thread
	writerThread *Thread

	ReaderThreadThrashingInterval time.Duration
	WriterThreadThrashingInterval time.Duration

	messageReceiver MessageReceiver
	incoming        MemoryBlockQueue
	outgoing        MemoryBlockQueue
}

func DefaultClientReaderThreadThrashingInterval() time.Duration {
	return 100 * time.Millisecond
}

func DefaultClientWriterThreadThrashingInterval() time.Duration {
	return 10 * time.Millisecond
}

func (this *Client) Push(memoryBlock []byte) {
	this.outgoing.Push(memoryBlock)
}

func (this *Client) Pop() []byte {
	return this.incoming.Pop()
}

func CheckIfNetTimeoutError(e error) bool {
	var result = false
	var netError, castResult = e.(net.Error)
	if castResult {
		result = netError.Timeout()
	} else {
		result = false
	}
	return result
}

func (this *Client) readerThreadRoutine(thread *Thread) {
	var buffer = make([]byte, 1, 1)
	var connection net.Conn
	var dropConnection = func() {
		this.connection = nil
		connection.Close()
		connection = nil
	}
	var tryExtractMessage = func() {
		var incomingMessage = this.messageReceiver.Extract()
		if incomingMessage != nil {
			this.incoming.Push(incomingMessage)
		}
	}
	var readForward = func() bool {
		connection.SetReadDeadline(time.Now().Add(1000 * time.Millisecond))
		var readLength, readResult = connection.Read(buffer)
		if readResult == nil {
			if readLength == 1 {
				this.messageReceiver.Write(buffer[0])
				tryExtractMessage()
			}
		} else if false == CheckIfNetTimeoutError(readResult) {
			dropConnection()
		}
		var dataReceived = (connection != nil) && (readResult == nil) && (readLength == 1)
		return dataReceived
	}
	for thread.Active {
		connection = this.connection
		if connection != nil {
			for readForward() {
			}
		}
		time.Sleep(this.ReaderThreadThrashingInterval)
	}
}

func (this *Client) writerThreadRoutine(thread *Thread) {
	var connection net.Conn
	var dropConnection = func() {
		this.connection = nil
		connection.Close()
		connection = nil
	}
	var writeData = func(data []byte) bool {
		var result = false
		if connection != nil {
			var writeLength, writeResult = this.connection.Write(data)
			if (writeResult == nil) && (writeLength == len(data)) {
				result = true
			} else {
				dropConnection()
			}
		}
		return result
	}
	var writeMessageSize = func(message []byte) bool {
		var messageSize = int64(len(message))
		var messageSizeMemoryBlock = Int64ToMemoryBlock(messageSize)
		var messageSizeData = messageSizeMemoryBlock.Bytes()
		return writeData(messageSizeData)
	}
	var writeMessage = func(message []byte) bool {
		var result = false
		if writeMessageSize(message) {
			if writeData(message) {
				result = true
			}
		}
		return result
	}
	var writeForward = func() bool {
		var outgoingMessage = this.outgoing.Pop()
		var writeResult = false
		if outgoingMessage != nil {
			writeResult = writeMessage(outgoingMessage)
		}
		return writeResult
	}
	for thread.Active {
		connection = this.connection
		if connection != nil {
			for writeForward() {
			}
		}
		time.Sleep(this.WriterThreadThrashingInterval)
	}
}

func (this *Client) Start() {
	this.readerThread = StartThread(this.readerThreadRoutine)
	this.writerThread = StartThread(this.writerThreadRoutine)
}

func (this *Client) Stop() {
	if this.writerThread != nil {
		this.writerThread.Active = false
		this.writerThread.WaitFor()
		this.writerThread = nil
	}
	if this.readerThread != nil {
		this.readerThread.Active = false
		this.readerThread.WaitFor()
		this.readerThread = nil
	}
	this.connection = nil
}
