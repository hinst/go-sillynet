package h_sillynet

import "bytes"

type MessageReceiver struct {
	memory                   bytes.Buffer
	expectedSize             int64
	expectedSizeData         Int64MemoryBlock
	expectedSizeDataPosition byte
}

func (this *MessageReceiver) WriteSize(x []byte) []byte {
	for (len(x) > 0) && (this.expectedSizeDataPosition < SizeOfInt64) {
		this.expectedSizeData[this.expectedSizeDataPosition] = x[0]
		this.expectedSizeDataPosition++
		x = x[1:]
		if this.expectedSizeDataPosition == SizeOfInt64 {
			this.expectedSize = MemoryBlockToInt64(this.expectedSizeData)
		}
	}
	return x
}

// x must be not nil.
func (this *MessageReceiver) Write(x []byte) {
	if this.expectedSizeDataPosition < SizeOfInt64 {
		x = this.WriteSize(x)
	}
	if len(x) > 0 {
		this.memory.Write(x)
	}
}

func (this *MessageReceiver) ready() bool {
	return (this.expectedSizeDataPosition == SizeOfInt64) && (this.expectedSize <= int64(this.memory.Len()))
}

func (this *MessageReceiver) clear() {
	this.memory.Reset()
	this.expectedSizeDataPosition = 0
}

func (this *MessageReceiver) Extract() []byte {
	if this.ready() {
		var result = this.memory.Bytes()
		if nil == result {
			result = make([]byte, 0)
		}
		this.clear()
		if this.expectedSize < int64(len(result)) {
			this.Write(result[this.expectedSize:])
			result = result[:this.expectedSize]
		}
		return result
	} else {
		return nil
	}
}
