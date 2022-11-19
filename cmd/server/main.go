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

type serverStruct struct {
	clients []*clientStruct
	lock    *sync.Mutex
}

type clientStruct struct {
	user       string
	connection net.Conn
	lock       *sync.Mutex
}

func process(client *clientStruct, server *serverStruct) {
	connection := client.connection
	write(client, "Connected to server successfully\n")
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
		case "\\quit":
			delUser(client, server)
			break ServerIn
		case "\\all":
			broadcast(client, server, "\\b "+client.user+": "+msg+"\n")
			api.Log(client.user + " broadcasted a message")
			break
		case "\\dm":
			user, newMsg, _ := strings.Cut(msg, " ")
			directMessage(client, server, user, newMsg)
			break
		case "\\online":
			showOnline(client, server)
			break
		case "\\help":
			showHelp(client, server)
			break
		}
	}
	connection.Close()
	api.Log("Client \"" + client.user + "\" disconnected\n")
}

func showHelp(client *clientStruct, server *serverStruct) {
	commands := "\"\\all <message>\" to broadcast <message> to all online users\n"
	commands += "\"\\dm <username> <message>\" to send <message> to <username>\n"
	commands += "\"\\online\" to view all currently online users\n"
	commands += "\"\\quit\" to log off the server\n"
	write(client, "\\b "+commands)
}

func showOnline(client *clientStruct, server *serverStruct) {
	users := ""
	for _, sClient := range server.clients {
		if sClient == client {
			users += sClient.user + " (You)\n"
		} else {
			users += sClient.user + "\n"
		}
	}
	write(client, "\\s Online users:\n"+users)
}

func directMessage(client *clientStruct, server *serverStruct, user string, msg string) {
	for _, sClient := range server.clients {
		if sClient.user == user {
			write(sClient, "\\d "+user+": "+msg+"\n")
			api.Log(client.user + " direct messaged " + user)
			return
		}
	}
	write(client, "\\e User \""+user+"\" does not exist\n")
}

func broadcast(client *clientStruct, server *serverStruct, msg string) {
	for _, sClient := range server.clients {
		if sClient != client {
			write(sClient, msg)
		}
	}
}

func write(client *clientStruct, msg string) {
	client.lock.Lock()
	defer client.lock.Unlock()
	client.connection.Write([]byte(msg))
}

func addUser(client *clientStruct, server *serverStruct) bool {
	server.lock.Lock()
	defer server.lock.Unlock()
	for i := 0; i < len(server.clients); i++ {
		if server.clients[i].user == client.user {
			write(client, "decline\n")
			return false
		}
	}
	server.clients = append(server.clients, client)
	write(client, "accept\n")
	api.Log("User \"" + client.user + "\" added\n")
	return true
}

func delUser(client *clientStruct, server *serverStruct) {
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
	server, err := net.Listen(serverType, host+":"+port)
	if err != nil {
		api.Err("Error starting server: " + err.Error() + "\n")
		os.Exit(1)
	}
	defer server.Close()
	api.Stat("Server started on " + host + "/" + port + "\n")
	newServer := serverStruct{*new([]*clientStruct), new(sync.Mutex)}
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
		newClient := clientStruct{user, connection, new(sync.Mutex)}
		for !addUser(&newClient, &newServer) {
			len, _ = connection.Read(reader)
			input = string(reader[:len])
			_, user, _ := strings.Cut(input, " ")
			newClient.user = user
		}
		go process(&newClient, &newServer)
	}

}
