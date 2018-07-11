package main

import (
	"flag"
	"log"
	"net/url"

	"github.com/go-chat-client/input"

	"github.com/go-chat-client/message"
	"github.com/gorilla/websocket"
)

func main() {
	var addr = flag.String("addr", "localhost:8000", "http service address")

	flag.Parse()
	log.SetFlags(0)

	url := url.URL{Scheme: "ws", Host: *addr, Path: "/echo"}
	log.Printf("connecting to %s", url.String())

	conn, _, err := websocket.DefaultDialer.Dial(url.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer conn.Close()
	done := make(chan struct{})

	go message.Read(conn, done)

	email := input.Scan("Email: ")
	user := input.Scan("Username: ")

	message.SendLoginRequest(conn, email, user)

	err = message.Send(conn, &done, email, user)
	if err != nil {
		log.Println("Send error:", err)
	}
}
