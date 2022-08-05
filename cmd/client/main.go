package main

import (
	"bufio"
	"net"
	"os"
	"strings"
	"sync"
	"tcp-chat/api"
)

const (
	serverType = "tcp"
)

func write(connection net.Conn, wg *sync.WaitGroup) {
	api.Broadcast("Enter commands:\n")
	defer wg.Done()
	for {
		reader := bufio.NewReader(os.Stdin)
		bLine, _, _ := reader.ReadLine()
		msg := string(bLine)
		_, err := connection.Write([]byte(msg))
		if err != nil {
			api.Err("Error sending message\n")
		}
		if msg == "\\quit" {
			break
		}
	}
	api.Stat("Sent disconnect signal to server " + connection.RemoteAddr().String() + "\n")
}

func read(connection net.Conn, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		reader := make([]byte, 10240)
		len, _ := connection.Read(reader)
		input := string(reader[:len])
		if input == "\\dc" {
			break
		}
		if input == "" {
			api.Err("Connection lost to server\n")
			break
		}
		sender, msg, _ := strings.Cut(input, " ")
		switch sender {
		case "\\b":
			api.Broadcast(msg)
			break
		case "\\e":
			api.Err(msg)
			break
		case "\\s":
			api.Stat(msg)
		case "\\d":
			api.DirectMessage(msg)
		default:
			api.Print(msg)
		}
	}
	connection.Close()
	api.Stat("Disconnected from server " + connection.RemoteAddr().String() + "\n")
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	api.Print("Enter server address: ")
	bHost, _, _ := reader.ReadLine()
	host := string(bHost)
	api.Print("Enter server port: ")
	bPort, _, _ := reader.ReadLine()
	port := string(bPort)
	if !api.CheckHost(host) {
		api.Err("Invalid address entered\n")
		os.Exit(1)
	} else if !api.CheckPort(port) {
		api.Err("Invalid port entered\n")
		os.Exit(1)
	}
	//establish connection
	connection, err := net.Dial(serverType, host+":"+port)
	if err != nil {
		api.Err("Cannot connect to server. Are you sure the server is running?\n")
		os.Exit(1)
	}
	for {
		api.Print("Enter username: ")
		user, _, _ := reader.ReadLine()
		if strings.ContainsAny(string(user), " \\") {
			api.Err("Username cannot contains spaces or \"\\\" characters\n")
			continue
		}
		_, err := connection.Write([]byte("\\user " + string(user)))
		if err != nil {
			api.Err("Error sending message\n")
			continue
		}

		reader := make([]byte, 10240)
		len, _ := connection.Read(reader)
		input := string(reader[:len])
		if input == "accept\n" {
			break
		}
		api.Err("Username already taken\n")
	}

	connected := make([]byte, 10240)
	length, _ := connection.Read(connected)
	msg := string(connected[:length])
	api.Stat(msg)

	var wg sync.WaitGroup
	wg.Add(2)
	go write(connection, &wg)
	go read(connection, &wg)
	wg.Wait()
}
