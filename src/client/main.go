package main

import (
	"client/utils"
	"fmt"
)

func main() {

	var ipAddr, userName string
	fmt.Print("Enter IP Address of server: ")
	fmt.Scanf("%s", &ipAddr)
	fmt.Print("Enter username: ")
	fmt.Scanf("%s", &userName)
	utils.StartClient(ipAddr, ":5000", userName)
	fmt.Println("Client completed")
}
