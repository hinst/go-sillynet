package main

import "os"
import "bufio"
import "fmt"
import "strings"
import "strconv"
import "time"
import "h_sillynet"

type Application struct {
	simpleServer   *h_sillynet.SimpleServer
	receiverThread *h_sillynet.Thread
}

func (application *Application) startReceiverThread() {
	application.receiverThread = h_sillynet.StartThread(func(thread *h_sillynet.Thread) {
		for thread.Active {
			time.Sleep(100 * time.Millisecond)
			var client = application.simpleServer.Client()
			if client != nil {
				var message = client.Pop()
				if message != nil {
					var messageText = string(message)
					fmt.Println("Message received: '" + messageText + "'")
				}
			}
		}
	})
}

func (application *Application) run() {
	application.simpleServer = &h_sillynet.SimpleServer{}
	application.simpleServer.Port = 9077
	application.startReceiverThread()
	var reader = bufio.NewReader(os.Stdin)
	var command = ""
	for command != "exit" {
		fmt.Print(">")
		command, _ = reader.ReadString('\n')
		command = strings.TrimSpace(command)
		if command == "start" {
			fmt.Println("Now starting server at port " + strconv.Itoa(application.simpleServer.Port) + "...")
			var startResult = application.simpleServer.Start()
			fmt.Println("Start result = " + strconv.FormatBool(startResult))
		} else if command == "stop" {
			fmt.Println("Now stopping server...")
			application.simpleServer.Stop()
		} else if strings.Index(command, "'") >= 0 {
			var client = application.simpleServer.Client()
			if client != nil {
				var messageText = command[1:]
				var messageData = []byte(messageText)
				client.Push(messageData)
			} else {
				fmt.Println("Can not send message: client not present.")
			}
		}
	}
}

func main() {
	var application = &Application{}
	application.run()
}
