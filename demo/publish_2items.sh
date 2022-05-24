#/!bin/bash

go run ~/Documents/go/src/github.com/nats-io/stan.go/examples/stan-pub/main.go orders "$(cat 2items.json)"

