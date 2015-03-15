package h_sillynet

import "net"

type Client struct {
	connection net.Conn
	incoming   MemoryBlockQueue
	outgoing   MemoryBlockQueue
}

func (this *Client) Push(memoryBlock []byte) {
	this.outgoing.Push(memoryBlock)
}

func (this *Client) Pop() []byte {
	return this.incoming.Pop()
}
