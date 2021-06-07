package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/marcelo-r/rabbitmq-example/cmd"
)

var command = flag.String("run", "", "which mode to run, options are: producer & consumer")
var delay int

func main() {
	flag.IntVar(&delay, "delay", 0, "force delay of sending to queue")

	flag.Parse()

	delayTime := time.Duration(delay) * time.Millisecond
	if *command == "producer" {
		cmd.Produce(delayTime)
	} else if *command == "consumer" {
		//cmd.Consumer()
		fmt.Println("producing")
	} else {
		fmt.Println("invalid option")
	}
}
