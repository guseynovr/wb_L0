package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os/user"

	stan "github.com/nats-io/stan.go"

	_ "github.com/lib/pq"
)

const (
	durName = "orders_L0"
)

var (
	db		*sql.DB
	sc		stan.Conn
	cache	map[string]Order
	tmpl	*template.Template
)

func init() {
	cache = make(map[string]Order)
	port := 5432
	u, err := user.Current()
	if err == nil && u.Username == "dgidget" {
		port = 5433 //for running at school
	}
	connStr := fmt.Sprintf(
		"postgres://wb_user:wb_pass@localhost:%d/wb?sslmode=disable", port)
	db, err = sql.Open("postgres", connStr)
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
	restoreCache()
	tmpl, err = template.ParseFiles("index.html")
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	defer db.Close()

	ins, err := db.Prepare("INSERT INTO orders (uid, data) values ($1, $2);")
	if err != nil {
		log.Fatal(err, "1")
	}
	defer ins.Close()
	
	upd, err := db.Prepare("UPDATE orders SET data = $2 WHERE uid = $1;")
	if err != nil {
		log.Fatal(err, "2")
	}
	defer upd.Close()

	defer sc.Close()
	sub, err := sc.Subscribe("orders", func (m *stan.Msg) {
		handleMessage(m, ins, upd)
	}, stan.DurableName(durName))
	defer sub.Unsubscribe()

	log.Fatal(http.ListenAndServe(":8000", http.HandlerFunc(handler)))
}

func handler(w http.ResponseWriter, r *http.Request) {
	uid := r.URL.Query().Get("uid")
	order := cache[uid]
	
	if err := tmpl.Execute(w, order); err != nil {
		log.Fatal(err)
	}
}

func restoreCache() {
	sel, err := db.Prepare("SELECT * FROM orders;")
	if err != nil {
		log.Fatal(err, "3")
	}
	defer sel.Close()
	
	rows, err := sel.Query()
	if err != nil {
		log.Fatal(err, "4")
	}
	var id int
	var uid, raw []byte
	var newOrder Order
	for rows.Next() {
		if err := rows.Scan(&id, &uid, &raw); err != nil {
			log.Fatal(err)
		}
		if err := json.Unmarshal(raw, &newOrder); err != nil {
			log.Fatal(err)
		}
		cache[newOrder.OrderUID] = newOrder
	}
}

func handleMessage(m *stan.Msg, ins, upd *sql.Stmt) {
	var newOrder Order;

	if err := json.Unmarshal(m.Data, &newOrder); err != nil {
		log.Println("Received invalid JSON data")
		return
	}

	_, exists := cache[newOrder.OrderUID]
	cache[newOrder.OrderUID] = newOrder;
	log.Printf("Received a new order: %s\n", newOrder.OrderUID)

	var op *sql.Stmt
	if exists {
		op = upd
	} else {
		op = ins
	}
	if _, err := op.Exec(newOrder.OrderUID, string(m.Data)); err != nil {
		log.Fatal(err)
	}
}


