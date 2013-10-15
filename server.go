package main

import (
	"fmt"
	"net/http"
//	"log"
	"time"
	"flag"
	
	"code.google.com/p/go.net/websocket"
)

const SERVICE = 8080
const HOST = ""
const FEED_INTV = 60 * time.Second

var p = Pool {
	connections: make(map[*Connection]bool),
	subscribe: make(chan *Connection),
	unsubscribe: make(chan *Connection),
}

func feed(conn *Connection) {
	n := 0
	for {
		message := fmt.Sprintf("Chunk: %d", n)
		err := websocket.Message.Send(conn.ws, message)
		if err != nil {
			break
		}
		time.Sleep(FEED_INTV)
		n++
	}
}

func socketHandler(ws *websocket.Conn) {
	var conn = &Connection {
		ws: ws,
	}
	
	p.subscribe <- conn
	feed(conn)
	p.unsubscribe <- conn
}

func displayStats() {
	for {
                fmt.Printf("open connections: %d\n", len(p.connections))
		time.Sleep(1 * time.Second)
        }
}

func main() {
	var host = flag.String("host", HOST, "Server listen host")
	var port = flag.Int("port", SERVICE, "Server listen port")
	
	flag.Parse()

	go p.Dispatch()
	go displayStats()
	
	http.Handle("/", websocket.Handler(socketHandler))
	err := http.ListenAndServe(fmt.Sprintf("%s:%d", *host, *port), nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
