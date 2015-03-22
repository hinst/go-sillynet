package h_sillynet

import "net"
import "time"
import "sync"

// "Simple" means that it has one access point. Only one client can connect.
type SimpleServer struct {
	Port                 int
	listener             *net.TCPListener
	acceptionThread      *Thread
	clientLocker         sync.Mutex
	client               *Client
	IncomingMessageEvent Event
}

func (this *SimpleServer) ClientAcceptionRoutine(thread *Thread) {
	var attachClient = func(connection net.Conn) {
		var client = NewClient()
		client.connection = connection
		client.IncomingMessageEvent = this.IncomingMessageEvent
		client.Start()
		this.SetClient(client)
	}
	var tryAcceptConnection = func() {
		this.listener.SetDeadline(time.Now().Add(1 * time.Second))
		var acceptedConnection, acceptResult = this.listener.Accept()
		if nil == acceptResult /*success*/ {
			attachClient(acceptedConnection)
		}
	}
	for thread.Active {
		if nil == this.Client() {
			tryAcceptConnection()
		} else {
			time.Sleep(1 * time.Second)
		}
	}
}

func (this *SimpleServer) Start() bool {
	var result = false
	if this.acceptionThread != nil {
		result = true // already started
	} else {
		var address = &net.TCPAddr{}
		address.Port = this.Port
		var listener, listenResult = net.ListenTCP("tcp", address)
		if nil == listenResult /*success*/ {
			this.listener = listener
			this.acceptionThread = StartThread(this.ClientAcceptionRoutine)
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
	var client = this.Client()
	this.SetClient(nil)
	if client != nil {
		client.Stop()
	}
}

func (this *SimpleServer) Client() *Client {
	this.clientLocker.Lock()
	var result = this.client
	this.clientLocker.Unlock()
	return result
}

func (this *SimpleServer) SetClient(a *Client) {
	this.clientLocker.Lock()
	this.client = a
	this.clientLocker.Unlock()
}
