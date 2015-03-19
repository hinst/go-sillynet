package main

import "os"
import "bufio"
import "h_sillynet"
import "fmt"
import "strings"
import "strconv"

func main() {
	var simpleServer h_sillynet.SimpleServer
	var receiverThread = h_sillynet.StartThread(func(thread *h_sillynet.Thread) {
		for thread.Active {
			var message = simpleServer.Client().Pop()
			if message != nil {
				var messageText = string(message)
				fmt.Println("Message received: '" + messageText + "'")
			}
		}
	})
	simpleServer.Port = 9077
	var reader = bufio.NewReader(os.Stdin)
	var command = ""
	for command != "exit" {
		fmt.Print(">")
		command, _ = reader.ReadString('\n')
		command = strings.TrimSpace(command)
		if command == "start" {
			fmt.Println("Now starting server at port " + strconv.Itoa(simpleServer.Port) + "...")
			simpleServer.Start()
		} else if command == "stop" {
			fmt.Println("Now stopping server...")
			simpleServer.Stop()
		}
	}
}
