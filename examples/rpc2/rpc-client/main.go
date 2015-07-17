package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

//	"gopkg.in/Miha-ha/turnpike.v2"
	"github.com/Miha-ha/turnpike"
)

var (
	client *turnpike.Client
)

func main() {
	turnpike.Debug()
	c, err := turnpike.NewWebsocketClient(turnpike.JSON, "ws://127.0.0.1:8000/ws")
	if err != nil {
		log.Fatal(err)
	}

	client = c

	_, err = c.JoinRealm("turnpike.examples", turnpike.ALLROLES, nil)
	if err != nil {
		log.Fatal(err)
	}

	quit := make(chan bool)
	c.Subscribe("alarm.ring", func([]interface{}, map[string]interface{}) {
		fmt.Println("The alarm rang!")
		quit <- true
	})

	if err := c.Register("alarm.set", alarmSet); err != nil {
		panic(err)
	}


	fmt.Print("Enter the timer duration: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	if err := scanner.Err(); err != nil {
		log.Fatalln("reading stdin:", err)
	}
	text := scanner.Text()
	if duration, err := strconv.Atoi(text); err != nil {
		log.Fatalln("invalid integer input:", err)
	} else {
		if _, err := c.Call("alarm.set", []interface{}{duration}, nil); err != nil {
			log.Fatalln("error setting alarm:", err)
		}
	}
	<-quit
}

// takes one argument, the (integer) number of seconds to set the alarm for
func alarmSet(args []interface{}, kwargs map[string]interface{}) (result *turnpike.CallResult) {
	duration, ok := args[0].(float64)
	if !ok {
		return &turnpike.CallResult{Err: turnpike.URI("rpc-example.invalid-argument")}
	}
	go func() {
		time.Sleep(time.Duration(duration) * time.Second)
		client.Publish("alarm.ring", nil, nil)
	}()
	return &turnpike.CallResult{Args: []interface{}{"hello"}}
}