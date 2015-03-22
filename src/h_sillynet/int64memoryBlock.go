package h_sillynet

type Int64MemoryBlock [8]byte

const SizeOfInt64 = 8

func Int64ToMemoryBlock(x int64) Int64MemoryBlock {
	var block Int64MemoryBlock
	var i uint
	for i = 0; i < SizeOfInt64; i++ {
		block[i] = byte((x >> (i * 8)) & 0xFF)
	}
	return block
}

func MemoryBlockToInt64(block Int64MemoryBlock) int64 {
	var x int64 = 0
	var i uint
	for i = 0; i < SizeOfInt64; i++ {
		var currentValue = int64(block[i])
		x = x | (currentValue << (i * 8))
	}
	return x
}

func (this Int64MemoryBlock) Bytes() []byte {
	var bytes = make([]byte, SizeOfInt64)
	for i := 0; i < SizeOfInt64; i++ {
		bytes[i] = this[i]
	}
	return bytes
}
