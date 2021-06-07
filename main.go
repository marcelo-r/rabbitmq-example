package main

import (
	"flag"
	"fmt"

	"github.com/marcelo-r/rabbitmq-example/cmd"
)

var command = flag.String("run", "", "which mode to run, options are: producer & consumer")

func main() {
	flag.Parse()
	if *command == "producer" {
		cmd.Producer()
	} else if *command == "consumer" {
		//cmd.Consumer()
		fmt.Println("producing")
	} else {
		fmt.Println("invalid option")
	}
}
