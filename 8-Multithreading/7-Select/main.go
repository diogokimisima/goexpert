package main

import (
	"fmt"
	"sync/atomic"
	"time"
)

type Message struct {
	id  int32
	msg string
}

func main() {
	c1 := make(chan Message)
	c2 := make(chan Message)
	var i int32 = 0

	// RabbitMQ
	go func() {
		for {
			atomic.AddInt32(&i, 1)
			msg := Message{i, "Hello from RabbitMQ"}
			time.Sleep(time.Second)
			c1 <- msg
		}
	}()

	// Kafka
	go func() {
		for {
			atomic.AddInt32(&i, 1)
			msg := Message{i, "Hello from Kafka"}
			time.Sleep(time.Second)
			c2 <- msg
		}
	}()

	for {
		select {
		case msg1 := <-c1: // rabbitmq
			fmt.Printf("Received ID messages %v from RabbitMQ: %s\n", msg1.id, msg1.msg)

		case msg2 := <-c2: // kafka
			fmt.Printf("Received ID messages %v from Kafka: %s\n", msg2.id, msg2.msg)

		case <-time.After(time.Second * 3):
			println("timeout")
		}
	}
}
