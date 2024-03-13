package main

import (
	"fmt"
	"time"

	"github.com/markbates/inflect"
	"github.com/shahzad73/GoUtils/test"
	"github.com/shahzad73/GoUtils/utils"
	"github.com/shahzad73/GoUtils/utils/Logs"
)

func main() {

	fmt.Println("Hello, Modules!")

	utils.PrintHello()

	logger := Logs.New(time.RFC3339, true)

	logger.Log("info", "starting up service")
	logger.Log("warning", "no tasks found")
	logger.Log("error", "exiting: no work performed")

	test.PrintHelloTest()

	ss := inflect.Camelize("This is some")
	fmt.Println(ss)
}
