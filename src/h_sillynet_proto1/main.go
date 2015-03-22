package main

import "os"
import "bufio"
import "fmt"
import "strings"
import "strconv"
import "time"
import "h_sillynet"

type Application struct {
	messageSendMoment time.Time
	simpleServer      *h_sillynet.SimpleServer
	receiverThread    *h_sillynet.Thread
}

func (this *Application) startReceiverThread() {
	this.receiverThread = h_sillynet.StartThread(func(thread *h_sillynet.Thread) {
		for thread.Active {
			var client = this.simpleServer.Client()
			if client != nil {
				var message = client.Pop()
				if message != nil {
					fmt.Println("ping ", time.Since(this.messageSendMoment))
					var messageText = ""
					if len(message) > 0 {
						messageText = string(message)
					}
					fmt.Println("Message received: '" + messageText + "'")
				}
			}
			time.Sleep(5 * time.Millisecond)
		}
	})
}

func (this *Application) stopReceiverThread() {
	this.receiverThread.Active = false
	this.receiverThread.WaitFor()
	this.receiverThread = nil
}

func (this *Application) run() {
	this.simpleServer = &h_sillynet.SimpleServer{}
	this.simpleServer.Port = 9077
	this.startReceiverThread()
	var reader = bufio.NewReader(os.Stdin)
	var command = ""
	for command != "exit" {
		fmt.Print(">")
		command, _ = reader.ReadString('\n')
		command = strings.TrimSpace(command)
		if command == "start" {
			fmt.Println("Now starting server at port " + strconv.Itoa(this.simpleServer.Port) + "...")
			var startResult = this.simpleServer.Start()
			fmt.Println("Start result = " + strconv.FormatBool(startResult))
		} else if command == "stop" {
			fmt.Println("Now stopping server...")
			this.simpleServer.Stop()
		} else if strings.Index(command, "'") >= 0 {
			var client = this.simpleServer.Client()
			if client != nil {
				var messageText = command[1:]
				fmt.Println("sending '" + messageText + "'")
				var messageData = []byte(messageText)
				client.Push(messageData)
				this.messageSendMoment = time.Now()
			} else {
				fmt.Println("Can not send message: client not present.")
			}
		}
	}
	this.stopReceiverThread()
}

func main() {
	var application = &Application{}
	application.run()
}
