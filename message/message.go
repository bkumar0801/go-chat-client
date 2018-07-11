package message

import (
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/go-chat-client/input"

	"github.com/gorilla/websocket"
)

/*
Message ...
*/
type Message struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Message  string `json:"message"`
}

/*
Read ...
*/
func Read(conn *websocket.Conn, done chan struct{}) {
	defer close(done)
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("reads:", err)
			return
		}
		var data map[string]interface{}
		if err := json.Unmarshal(message, &data); err != nil {
			log.Println("reads:", err)
			return
		}

		log.Printf("%s: %s\n", data["username"].(string), data["message"].(string))
	}
}

/*
Send ...
*/
func Send(conn *websocket.Conn, done chan struct{}, email, user string) error {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	msg := make(chan Message)
	defer close(msg)
	message := Message{
		Email:    email,
		Username: user,
	}
	log.Println("You can send messages now:> ")
	for {
		message.Message = input.Scan("")
		go func() {
			if message.Message != "" {
				msg <- message
			}
		}()

		select {
		case <-done:
			log.Println("Done")
			return nil
		case chat := <-msg:
			b, _ := json.Marshal(chat)
			err := conn.WriteMessage(websocket.BinaryMessage, b)
			if err != nil {
				log.Println("write:", err)
				return err
			}
		case <-interrupt:
			log.Println("interrupt")

			/* Cleanly close the connection by sending a close message and then
			waiting (with timeout) for the server to close the connection.
			*/
			err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return err
			}
			select {
			case <-done:
				log.Println("Done")
			case <-time.After(time.Second):
			}
			return nil
		}
	}
}

/*
SendLoginRequest ...
*/
func SendLoginRequest(conn *websocket.Conn, email, user string) {
	message := Message{
		Email:    email,
		Username: user,
		Message:  "login",
	}
	data, _ := json.Marshal(message)
	err := conn.WriteMessage(websocket.BinaryMessage, data)
	if err != nil {
		log.Println("write:", err)
		return
	}

}
