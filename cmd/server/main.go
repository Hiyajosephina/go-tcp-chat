// Package main contains all the functions utilized by the server
// application to both connect and communicate with all the clients
package main

import (
	"bufio"
	"net"
	"os"
	"strings"
	"sync"
	"tcp-chat/api"
)

// Server is a struct containing the server lock and array of clients
type Server struct {
	clients []*Client
	lock    *sync.Mutex
}

// Client is a struct containing data related to a individual clients
type Client struct {
	user       string
	connection net.Conn
	lock       *sync.Mutex
}

// Process facilitates communication with the server and an individual
// client. Reading and writing to the client's net.Conn as well
// as computing which command was received from the client
// and calling the appropriate function in order to process the command
func Process(client *Client, server *Server) {
	connection := client.connection
	Write(client, "Connected to server successfully\n")

// ServerIn is the main loop that for client-server communication.
// Exiting the loop indicates the client is no longer connected to 
// the server or that the quit command was called
ServerIn:
	for {
		reader := make([]byte, 10240)
		len, _ := connection.Read(reader)
		input := string(reader[:len])
		if input == "" {
			continue
		}
		cmd, msg, _ := strings.Cut(input, " ")
		switch cmd {
		
		// "\quit" is the quit command
		case "\\quit":
			DelUser(client, server)
			break ServerIn

		// "\all <message>" is the broadcast command
		case "\\all":
			Broadcast(client, server, "\\b "+client.user+": "+msg+"\n")
			break

		// "\dm <username> <message>" is the direct message command
		case "\\dm":
			user, newMsg, _ := strings.Cut(msg, " ")
			DirectMessage(client, server, user, newMsg)
			break

		// "\online" is the show online command
		case "\\online":
			ShowOnline(client, server)
			break

		// "\help" is the help command
		case "\\help":
			ShowHelp(client, server)
			break
		}
	}
	// close connection after exiting ServerIn loop and log the disconnection
	connection.Close()
	api.Log("Client \"" + client.user + "\" disconnected\n")
}

// ShowHelp writes to the client that they have to display the instructions 
// on how to issue a command to the server and in what format
func ShowHelp(client *Client, server *Server) {
	Write(client, "\\h ")
}

// ShowOnline writes to the client all the currently online users as well as
// indicating which user they are
func ShowOnline(client *Client, server *Server) {
	users := ""
	for _, sClient := range server.clients {
		if sClient == client {
			users += sClient.user + " (You)\n"
		} else {
			users += sClient.user + "\n"
		}
	}
	Write(client, "\\s Online users:\n"+users)
}

// DirectMessage lets a client write a direct message to another client by searching
// for that clients username and then writing to their net.Conn
func DirectMessage(client *Client, server *Server, user string, msg string) {
	for _, sClient := range server.clients {
		if sClient.user == user {
			Write(sClient, "\\d "+client.user+": "+msg+"\n")
			api.Log(client.user + " direct messaged " + user + "\n")
			return
		}
	}
	Write(client, "\\e User \""+user+"\" does not exist\n")
}

// Broadcast lets a client write a message that will be sent to all other clients
// that are connected to the server
func Broadcast(client *Client, server *Server, msg string) {
	for _, sClient := range server.clients {
		if sClient != client {
			Write(sClient, msg)
		}
	}
	api.Log(client.user + " broadcasted a message\n")
	
}

// Write writes to the client's net.Conn
func Write(client *Client, msg string) {
	client.lock.Lock()
	defer client.lock.Unlock()
	client.connection.Write([]byte(msg))
}

// AddUser adds a client to the client array struct if the username is unique
func AddUser(client *Client, server *Server) bool {
	server.lock.Lock()
	defer server.lock.Unlock()
	for i := 0; i < len(server.clients); i++ {
		if server.clients[i].user == client.user {
			Write(client, "decline\n")
			return false
		}
	}
	server.clients = append(server.clients, client)
	Write(client, "accept\n")
	api.Log("User \"" + client.user + "\" added\n")
	return true
}

// DelUser removes a client from the client array once they have disconnected
func DelUser(client *Client, server *Server) {
	server.lock.Lock()
	defer server.lock.Unlock()
	for i := 0; i < len(server.clients); i++ {
		if server.clients[i].user == client.user {
			server.clients = append(server.clients[:i], server.clients[i+1:]...)
			api.Log("User \"" + client.user + "\" removed\n")
			return
		}
	}
	api.Log("User " + client.user + " not found\n")
}

// main is the entry point for the server application. The initial connection
// between server and client occurs here, as well as connecting new clients
// and creating/modifying the appropriate structs for these new connections
func main() {
	reader := bufio.NewReader(os.Stdin)
	api.Print("Enter server address: ")
	bHost, _, _ := reader.ReadLine()
	host := string(bHost)
	api.Print("Enter server port: ")
	bPort, _, _ := reader.ReadLine()
	port := string(bPort)
	if !api.CheckHost(host) {
		api.Err("Invalid address entered")
		os.Exit(0)
	} else if !api.CheckPort(port) {
		api.Err("Invalid port entered")
		os.Exit(0)
	}
	api.Stat("Starting up server...\n")
	server, err := net.Listen("tcp", host+":"+port)
	if err != nil {
		api.Err("Error starting server: " + err.Error() + "\n")
		os.Exit(1)
	}
	defer server.Close()
	api.Stat("Server started on " + host + "/" + port + "\n")
	newServer := Server{*new([]*Client), new(sync.Mutex)}
	for {
		connection, err := server.Accept()
		if err != nil {
			api.Err("Error accepting: " + err.Error())
			os.Exit(1)
		}
		api.Log("Client " + connection.RemoteAddr().String() + " successfully connected\n")
		reader := make([]byte, 10240)
		len, _ := connection.Read(reader)
		input := string(reader[:len])
		_, user, _ := strings.Cut(input, " ")
		newClient := Client{user, connection, new(sync.Mutex)}
		for !AddUser(&newClient, &newServer) {
			len, _ = connection.Read(reader)
			input = string(reader[:len])
			_, user, _ := strings.Cut(input, " ")
			newClient.user = user
		}
		go Process(&newClient, &newServer)
	}

}
