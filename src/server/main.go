package main

import (
	"fmt"
	"server/utils"
)

func main() {
	utils.StartServer(":5000")
	fmt.Println("Server stopped")
}
