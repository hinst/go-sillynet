package h_sillynet

import "net"
import "time"
import "sync"

func NewClient() *Client {
	var result = &Client{}
	result.ReaderThreadThrashingInterval = DefaultClientReaderThreadThrashingInterval()
	result.WriterThreadThrashingInterval = DefaultClientWriterThreadThrashingInterval()
	result.writerThreadEvent = sync.NewCond(&sync.Mutex{})
	return result
}

type Client struct {
	connection       net.Conn
	connectionLocker sync.Mutex

	readerThread      *Thread
	writerThread      *Thread
	writerThreadEvent *sync.Cond

	ReaderThreadThrashingInterval time.Duration
	WriterThreadThrashingInterval time.Duration

	messageReceiver MessageReceiver
	incoming        MemoryBlockQueue
	outgoing        MemoryBlockQueue
}

func DefaultClientReaderThreadThrashingInterval() time.Duration {
	return 0 * time.Millisecond
}

func DefaultClientWriterThreadThrashingInterval() time.Duration {
	return 10 * time.Millisecond
}

const DefaultReadBufferSize = 16

func (this *Client) Push(memoryBlock []byte) {
	this.outgoing.Push(memoryBlock)
	this.pokeWriterThread()
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
	var buffer = make([]byte, DefaultReadBufferSize)
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
		var readLength, readResult = connection.Read(buffer)
		if readResult == nil {
			if readLength > 0 {
				this.messageReceiver.Write(buffer[:readLength])
				tryExtractMessage()
			}
		} else if false == CheckIfNetTimeoutError(readResult) {
			dropConnection()
		}
		var dataReceived = (connection != nil) && (readResult == nil) && (readLength > 0)
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
		this.writerThreadEvent.L.Lock()
		this.writerThreadEvent.Wait()
		this.writerThreadEvent.L.Unlock()
		//time.Sleep(this.WriterThreadThrashingInterval)
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
		this.pokeWriterThread()
		this.readerThread.WaitFor()
		this.readerThread = nil
	}
	this.connection = nil
}

func (this *Client) pokeWriterThread() {
	this.writerThreadEvent.L.Lock()
	this.writerThreadEvent.Signal()
	this.writerThreadEvent.L.Unlock()
}
