package h_sillynet

import "net"
import "time"
import "sync"

type Thread struct {
	waitGroup sync.WaitGroup
	Active    bool
}

func StartThread(f func(thread *Thread)) *Thread {
	var thread = &Thread{}
	thread.waitGroup.Add(1)
	go func() {
		f(thread)
		thread.waitGroup.Done()
	}()
	return thread
}

type Client struct {
	connection net.Conn
}

// "Simple" means that it has a single access point. Only one client can connect.
type SimpleServer struct {
	Port            int
	listener        *net.TCPListener
	acceptionThread *Thread
	Active          bool
	client          *Client
}

func (simpleServer *SimpleServer) ClientAcceptionRoutine(thread *Thread) {
	var tryAcceptConnection = func() {
		simpleServer.listener.SetDeadline(time.Now().Add(1 * time.Second))
		var acceptedConnection, acceptResult = simpleServer.listener.Accept()
		if nil == acceptResult /*success*/ {
			simpleServer.client = &Client{connection: acceptedConnection}
		}
	}
	for simpleServer.Active {
		if nil == simpleServer.client {
			tryAcceptConnection()
		}
	}
}

func (simpleServer *SimpleServer) Start() bool {
	var result = false
	if simpleServer.Active {
		result = true // already started
	} else {
		var address = &net.TCPAddr{}
		address.Port = simpleServer.Port
		var listener, listenResult = net.ListenTCP("tcp", address)
		if nil == listenResult /*success*/ {
			var simpleServer = &SimpleServer{}
			simpleServer.listener = listener
			simpleServer.Active = true
			simpleServer.acceptionThread = StartThread(simpleServer.ClientAcceptionRoutine)
			result = true
		}
	}
	return result
}
