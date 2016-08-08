package utils

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"
)

// Map of available users
var Users []string

// IP Address of server
var IPAddr string

type Message struct {
	Cmd    string // PRIVATE or BROADCAST or USERLIST
	Client string // Sender
	Msg    string // Message
}

// Convert the string into Message struct
func ParseMessage(msg string) Message {
	var message Message
	msgBytes := []byte(msg)
	err := json.Unmarshal(msgBytes, &message)
	handleErr(err)
	fmt.Println("message:", message)
	//message.ExecuteCommand()
	return message
}

// Function to send message to the client
func (msg Message) ExecuteCommand() bool {

	if msg.Cmd == "USERLIST" {
		fmt.Println("Available users", msg.Client)
	} else if msg.Cmd == "BROADCAST" {
		fmt.Println("Broadcast message from", msg.Client, " : ", msg.Msg)
	} else if msg.Cmd == "PRIVATE" {
		fmt.Println("Private message from", msg.Client, " : ", msg.Msg)
	}
	return true
}

// Function to start the client
func StartClient(serverIP, portNo, userName string) {
	if joinServer(serverIP, portNo, userName) {
		IPAddr = serverIP
		fmt.Println("Available Users:", Users)
		var isStarted chan bool
		go startReceiver(isStarted, portNo)
		if <-isStarted {
			var s string
			for {
				fmt.Scanf("%s", &s)
				msgstr := strings.Split(s, ":")
				var message Message

				switch msgstr[0] {
				case "BROADCAST":
					message = Message{"BROADCAST", "", msgstr[1]}
					go sendMessage(message)
					break

				case "LEAVE":
					message = Message{"LEAVE", "", ""}
					go sendMessage(message)
					break

				case "SEND":
					message = Message{msgstr[0], msgstr[1], msgstr[2]}
					go sendMessage(message)
					break

				default:
					fmt.Println("Invalid command")

				}
				if msgstr[0] == "BROADCAST" {
					message = Message{"BROADCAST", "", msgstr[1]}
				} else if msgstr[0] == "LEAVE" {
					message = Message{"LEAVE", "", ""}
				} else {
					message = Message{"SEND", msgstr[0], msgstr[1]}
				}
				go sendMessage(message)
			}
		}
	}
}

// Function to join server with username
func joinServer(serverIP, portNo, userName string) bool {
	msg := Message{"JOIN", userName, userName}
	var responseMsg Message
	buf, _ := json.Marshal(msg)

	conn, err := net.Dial("tcp", serverIP+portNo)
	defer conn.Close()
	handleErr(err)

	_, err2 := conn.Write(buf)
	handleErr(err2)

	reply := make([]byte, 1024)

	n, err3 := conn.Read(reply)
	handleErr(err3)

	err4 := json.Unmarshal(reply[0:n], &responseMsg)
	if err4 != nil {
		handleErr(err4)
	}

	if responseMsg.Cmd == "USERLIST" {
		Users = strings.Split(responseMsg.Client, "|")
	}

	fmt.Println("Joined server successfully")
	return true
}

// Starts the receviver
func startReceiver(isStared chan bool, receiverPortNo string) {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", receiverPortNo)
	handleErr(err)

	listener, err := net.ListenTCP("tcp", tcpAddr)
	handleErr(err)
	fmt.Println("Client started !!")
	isStared <- true
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}

		go HandleReceiver(conn)
	}
}

// Recieves message from the client. Parse the message and call sendMessage method.
func HandleReceiver(conn net.Conn) {
	defer conn.Close()
	fmt.Println("Inside client handle receiver")
	var buf [1024]byte
	for {
		// read upto 1024 bytes
		n, err := conn.Read(buf[0:])

		if err != nil {
			return
		}
		fmt.Println("Inside handle receiver", len(buf), buf)
		msg := string(buf[0:n])
		fmt.Println("Inside client handle receiver", msg)
		message := ParseMessage(msg)
		message.ExecuteCommand()

	}

}

// Sends message to specified client
func sendMessage(msg Message) {
	fmt.Println("Inside SendMessage")
	fmt.Println(msg.Msg, msg.Client)

	conn, err := net.Dial("tcp", IPAddr+":5000")
	defer conn.Close()
	handleErr(err)

	buf := []byte(msg.Msg)
	_, err2 := conn.Write(buf)
	if err2 != nil {
		return
	}

}

func handleErr(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}
