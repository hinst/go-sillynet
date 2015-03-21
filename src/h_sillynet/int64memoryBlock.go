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
		x = x + int64(block[i]<<(i*8))
	}
	return x
}
