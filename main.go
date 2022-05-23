package main

import (
	"database/sql"
	"fmt"
	"log"
	"os/user"

	stan "github.com/nats-io/stan.go"

	_ "github.com/lib/pq"
)

var (
	db *sql.DB
	sc stan.Conn
)

func init() {
	port := 5432
	u, err := user.Current()
	if err != nil && u.Username == "dgidget" {
		port = 5433 //for running at school
	}
	connStr := fmt.Sprintf(
		"postgres://wb_user:wb_pass@localhost:%d/wb?sslmode=disable", port)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}
	sc, err = stan.Connect("test-cluster", "L0")
	if err != nil {
		log.Fatal(err)
	}
}

func main() {

}
