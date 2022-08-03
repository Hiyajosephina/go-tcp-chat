package main

import (
	"bufio"
	"fmt"
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
	fmt.Println("Enter commands:")
	defer wg.Done()
	for {
		reader := bufio.NewReader(os.Stdin)
		bLine, _, _ := reader.ReadLine()
		msg := string(bLine)
		_, err := connection.Write([]byte(msg))
		if err != nil {
			fmt.Println("Error sending message")
		}
		if msg == "\\quit" {
			break
		}
	}
	fmt.Printf("Sent disconnect signal to server %s\n", connection.RemoteAddr().String())
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
			fmt.Println("Connection lost to server")
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
		default:
			fmt.Print(msg)
		}
	}
	connection.Close()
	fmt.Printf("Disconnected from server %s\n", connection.RemoteAddr().String())
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter server address: ")
	bHost, _, _ := reader.ReadLine()
	host := string(bHost)
	fmt.Print("Enter server port: ")
	bPort, _, _ := reader.ReadLine()
	port := string(bPort)
	if !api.CheckHost(host) {
		fmt.Println("Invalid address entered")
		os.Exit(1)
	} else if !api.CheckPort(port) {
		fmt.Println("Invalid port entered")
		os.Exit(1)
	}
	//establish connection
	connection, err := net.Dial(serverType, host+":"+port)
	if err != nil {
		fmt.Println("Cannot connect to server. Are you sure the server is running?")
		os.Exit(1)
	}
	for {
		fmt.Print("Enter username: ")
		user, _, _ := reader.ReadLine()
		_, err := connection.Write([]byte("\\user " + string(user)))
		if err != nil {
			fmt.Println("Error sending message")
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

	var wg sync.WaitGroup
	wg.Add(2)
	go write(connection, &wg)
	go read(connection, &wg)
	wg.Wait()
}
