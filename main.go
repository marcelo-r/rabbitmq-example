package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/marcelo-r/rabbitmq-example/cmd"
)

var command = flag.String("run", "", "which mode to run, options are: produce & consume")
var filename = flag.String("f", "mock_data.csv", "file with mock data to enqueue in rabbitmq")
var delay int

func main() {
	flag.IntVar(&delay, "delay", 0, "force delay of sending to queue")

	flag.Parse()

	delayTime := time.Duration(delay) * time.Millisecond
	if *command == "producer" {
		_ = cmd.Produce(delayTime, *filename)
	} else if *command == "consumer" {
		cmd.Consume()
	} else {
		fmt.Println("invalid option")
	}
}
