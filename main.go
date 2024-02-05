package main

import (
	"fmt"
	"time"

	"github.com/shahzad73/GoUtils/utils"
	"github.com/shahzad73/GoUtils/utils/Logs"
)

func main() {

	fmt.Println("Hello, Modules!")

	utils.PrintHello()

	logger := Logs.New(time.RFC3339, true)

	logger.Log("This is a debug statement...")

}
