package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"
)

var Users map[string]string

type Message struct {
	Cmd    string
	Client string
	Msg    string
}

// Command executor
func (msg Message) ExecuteCommand(conn net.Conn) bool {

	switch msg.Cmd {
	case "JOIN":
		if Users == nil {
			fmt.Println("inside user")
			Users = make(map[string]string)
		}

		Users[msg.Msg] = strings.Split(conn.RemoteAddr().String(), ":")[0]
		msg = Message{"USERLIST", getUserList(), "SUCCESS"}
		buf, _ := json.Marshal(msg)
		conn.Write(buf)
		break
	case "BROADCAST":
		sender := strings.Split(conn.RemoteAddr().String(), ":")[0]
		for _, val := range Users {
			go sendMessage(msg.Msg, sender, val, "BROADCAST")
		}
		break
	case "SEND":
		sender := strings.Split(conn.RemoteAddr().String(), ":")[0]
		receiver := msg.Client
		go sendMessage(msg.Msg, sender, receiver, "SEND")
		break
	case "LEAVE":
		clientAddr := strings.Split(conn.RemoteAddr().String(), ":")[0]
		go removeClient(clientAddr)
		break
	default:
		fmt.Println("Not listed command")
	}

	return true
}

// Convert string message to Message struct
func ParseMessage(msg string) Message {
	var message Message
	msgBytes := []byte(msg)
	err := json.Unmarshal(msgBytes, &message)
	handleErr(err)
	fmt.Println("Person", message)
	return message
}

// Recieves message from the client. Parse the message and call sendMessage method.
func HandleReceiver(conn net.Conn) {
	defer conn.Close()
	var buf [1024]byte
	for {
		// read upto 1024 bytes
		n, err := conn.Read(buf[0:])

		if err != nil {
			return
		}
		msg := string(buf[0:n])
		fmt.Println("Messaged received from client: ", msg)
		message := ParseMessage(msg)
		message.ExecuteCommand(conn)
		break
	}
	fmt.Println("Done handle Receiver", Users)

}

// Sends message to receiver
func sendMessage(msg, sender, receiver, msgType string) {
	fmt.Println("Inside send message :", msg, sender, receiver)
	conn, err := net.Dial("tcp", Users[receiver]+":5000")
	defer conn.Close()
	handleErr(err)
	message := Message{msgType, sender, msg}
	buf, err2 := json.Marshal(message)
	handleErr(err2)
	_, err3 := conn.Write(buf)
	if err3 != nil {
		return
	}
	fmt.Println("Done SendMessage")
}

// Remove client from the list
func removeClient(clientAddr string) {
	for key, value := range Users {
		if value == clientAddr {
			delete(Users, key)
			conn, err := net.Dial("tcp", value+":5000")
			defer conn.Close()
			handleErr(err)
			msg := Message{"OK", "", "SUCCESS"}
			buf, err1 := json.Marshal(msg)
			handleErr(err1)
			_, err3 := conn.Write(buf)
			if err3 != nil {
				return
			}
			break
		}
	}
}

// Retrieves the user list as "|" seperated value
func getUserList() string {
	var userlist string
	fmt.Println(len(Users))
	for key, value := range Users {
		fmt.Println("key", key, "value", value)
		userlist = userlist + key + "|"

	}
	return strings.TrimRight(userlist, "|")
}

// Start the server at specified portno
func StartServer(portNo string) {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", portNo)
	handleErr(err)

	listener, err := net.ListenTCP("tcp", tcpAddr)
	handleErr(err)
	fmt.Println("Server started successfully..!!")
	fmt.Println("")
	fmt.Println("")
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}

		go HandleReceiver(conn)
	}
}

func handleErr(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}
