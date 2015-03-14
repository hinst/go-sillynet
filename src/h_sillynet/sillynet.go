package h_sillynet

import "net"
import "strconv"
import "time"
import "sync"

type Thread struct {
	waitGroup sync.WaitGroup
	Active    bool
}

func StartThread(f func(thread *Thread)) Thread {
	var thread Thread
	thread.waitGroup.Add(1)
	go func() {
		f(thread)
		thread.waitGroup.Done()
	}()
}

// "Simple" means that it has a single access point. Only one client can connect.
type SimpleServer struct {
	Port            int
	listener        net.TCPListener
	acceptionThread sync.WaitGroup
	Active          bool
}

func (simpleServer *SimpleServer) ClientAcceptionRoutine(thread *Thread) {
	for simpleServer.Active {
		simpleServer.listener.SetDeadline(time.Now().Add(1 * time.Second))
		var acceptedConnection, acceptResult = simpleServer.listener.Accept()
		if nil == acceptResult /*success*/ {
			break
		}
	}
}

func (simpleServer *SimpleServer) Start() bool {
	if false == Active {
	}
	var result = false
	var listener, listenResult = net.Listen("tcp", ":"+strconv.Itoa(port))
	if nil == listenResult /*success*/ {
		var simpleServer = &SimpleServer{}
		simpleServer.listener = listener
		simpleServer.Active = true
		simpleServer.acceptionThread = StartThread(simpleServer.ClientAcceptionRoutine)
		result = true
	}
	return result
}
