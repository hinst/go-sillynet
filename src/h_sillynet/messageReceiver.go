package h_sillynet

import "bytes"

type MessageSizeData [8]byte

type MessageReceiver struct {
	sizeData MessageSizeData
	buffer   bytes.Buffer
	size     int64
}
