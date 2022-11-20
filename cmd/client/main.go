// Package main contains all the functions utilized by the client
// application to both connect and communicate with the server
package main

import (
	"bufio"
	"net"
	"os"
	"strings"
	"sync"
	"tcp-chat/api"
)

// Write writes to the net.Conn that connects the current client to the server.
// If the client sends the "\quit" command, the net.Conn is closed and the
// client is disconnected
func Write(connection net.Conn, wg *sync.WaitGroup) {
	api.Broadcast("Enter commands:\t\t\t(Enter \"\\help\" for list of commands)\n")
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

// Read reads from the net.Conn that connects the current client to the server.
// If the client reads the "\dc" command or a blank string, the net.Conn is
// closed and the client is disconnected
func Read(connection net.Conn, wg *sync.WaitGroup) {
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
		case "\\h":
			commands := "\"\\all <message>\" to broadcast <message> to all online users\n"
			commands += "\"\\dm <username> <message>\" to send <message> to <username>\n"
			commands += "\"\\online\" to view all currently online users\n"
			commands += "\"\\quit\" to log off the server\n"
			api.Broadcast(commands)
			break
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

// main is the entry point for the client application. The initial connection
// between server and client occurs here, as well as providing the clients
// preffered username
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
	connection, err := net.Dial("tcp", host+":"+port)
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
	go Write(connection, &wg)
	go Read(connection, &wg)
	wg.Wait()

}
