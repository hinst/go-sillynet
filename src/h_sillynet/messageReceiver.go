package h_sillynet

import "bytes"

type MessageReceiver struct {
	memory                   bytes.Buffer
	expectedSize             int64
	expectedSizeData         Int64MemoryBlock
	expectedSizeDataPosition byte
}

func (this *MessageReceiver) Write(x byte) {
	if this.expectedSizeDataPosition < SizeOfInt64 {
		this.expectedSizeData[this.expectedSizeDataPosition] = x
		this.expectedSizeDataPosition++
		if this.expectedSizeDataPosition == SizeOfInt64 {
			this.expectedSize = MemoryBlockToInt64(this.expectedSizeData)
		}
	} else {
		this.memory.WriteByte(x)
	}
}

func (this *MessageReceiver) ready() bool {
	return (this.expectedSizeDataPosition == SizeOfInt64) && (this.expectedSize == int64(this.memory.Len()))
}

func (this *MessageReceiver) clear() {
	this.memory.Reset()
	this.expectedSizeDataPosition = 0
}

func (this *MessageReceiver) Extract() []byte {
	if this.ready() {
		var result = this.memory.Bytes()
		this.clear()
		return result
	} else {
		return nil
	}
}
