package h_sillynet

import "net"
import "time"

// "Simple" means that it has a single access point. Only one client can connect.
type SimpleServer struct {
	Port            int
	listener        *net.TCPListener
	acceptionThread *Thread
	client          *Client
	incoming        *MemoryBlockQueue
	outgoing        *MemoryBlockQueue
}

func (simpleServer *SimpleServer) ClientAcceptionRoutine(thread *Thread) {
	var tryAcceptConnection = func() {
		simpleServer.listener.SetDeadline(time.Now().Add(1 * time.Second))
		var acceptedConnection, acceptResult = simpleServer.listener.Accept()
		if nil == acceptResult /*success*/ {
			simpleServer.client = &Client{connection: acceptedConnection}
		}
	}
	for thread.Active {
		if nil == simpleServer.client {
			tryAcceptConnection()
		}
	}
}

func (simpleServer *SimpleServer) Start() bool {
	var result = false
	if nil == simpleServer.acceptionThread {
		result = true // already started
	} else {
		var address = &net.TCPAddr{}
		address.Port = simpleServer.Port
		var listener, listenResult = net.ListenTCP("tcp", address)
		if nil == listenResult /*success*/ {
			var simpleServer = &SimpleServer{}
			simpleServer.listener = listener
			simpleServer.acceptionThread = StartThread(simpleServer.ClientAcceptionRoutine)
			result = true
		}
	}
	return result
}

func (this *SimpleServer) Stop() {
	if this.acceptionThread != nil {
		this.acceptionThread.Active = false
		this.acceptionThread.WaitFor()
		this.acceptionThread = nil
	}
}
