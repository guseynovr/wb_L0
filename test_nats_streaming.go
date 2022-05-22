package main

import (
	"fmt"
	"time"

	stan "github.com/nats-io/stan.go"
)

func main() {
	stanConn, err := stan.Connect("test-cluster", "test_user2")
	if err != nil {
		panic(err)
	}
	stanConn.Subscribe("orders", func(m *stan.Msg) {
		fmt.Printf("Received a message(test_user2, orders): %s\n",
			string(m.Data))
	})
	time.Sleep(10 * time.Second)
}
