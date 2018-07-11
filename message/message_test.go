package message

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

func TestSend(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer c.Close()
		for {
			mt, message, err := c.ReadMessage()
			if err != nil {
				break
			}
			err = c.WriteMessage(mt, message)
			if err != nil {
				break
			}
		}
	}))
	defer testServer.Close()

	// Convert http://127.0.0.1 to ws://127.0.0.1/echo
	u := "ws" + strings.TrimPrefix(testServer.URL, "http") + "/echo"

	// Connect to the server
	ws, _, err := websocket.DefaultDialer.Dial(u, nil)
	if err != nil {
		t.Fatalf("%v", err)
	}
	defer ws.Close()

	done := make(chan struct{})

	tmpfile, originalStdin := mockStdInput("Hello,Test")
	os.Stdin = tmpfile

	//Read back message received from server
	go func() {
		_, data, err := ws.ReadMessage()
		if err != nil {
			t.Fatalf("%v", err)
		}
		got := Message{}
		if err := json.Unmarshal(data, &got); err != nil {
			t.Fatalf("%v", err)
		}
		want := Message{
			Email:    "xyz@xyz.com",
			Username: "Test",
			Message:  "Hello,Test",
		}
		fmt.Printf("got: %v", got)
		if !reflect.DeepEqual(want, got) {
			t.Errorf("\nUnexpected response from server: \n\t\t expected: %q, \n\t\t actual: %q", want, got)
		}
	}()

	go func() {
		time.Sleep(2 * time.Millisecond)
		close(done)
	}()
	//Send message to server
	err = Send(ws, done, "xyz@xyz.com", "Test")
	if err != nil {
		t.Fatalf("%v", err)
	}
	os.Stdin = originalStdin // Restore original Stdin

	if err := tmpfile.Close(); err != nil {
		log.Fatal(err)
	}

}
func TestSendLoginRequest(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer c.Close()
		for {
			mt, message, err := c.ReadMessage()
			if err != nil {
				break
			}
			err = c.WriteMessage(mt, message)
			if err != nil {
				break
			}
		}
	}))
	defer testServer.Close()

	// Convert http://127.0.0.1 to ws://127.0.0.1/echo
	u := "ws" + strings.TrimPrefix(testServer.URL, "http") + "/echo"

	// Connect to the server
	ws, _, err := websocket.DefaultDialer.Dial(u, nil)
	if err != nil {
		t.Fatalf("%v", err)
	}
	defer ws.Close()
	SendLoginRequest(ws, "xyz@xyz.com", "Test")

	//Read back message received from server
	_, data, err := ws.ReadMessage()
	if err != nil {
		t.Fatalf("%v", err)
	}
	got := Message{}
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("%v", err)
	}
	want := Message{
		Email:    "xyz@xyz.com",
		Username: "Test",
		Message:  "login",
	}
	if !reflect.DeepEqual(want, got) {
		t.Errorf("\nUnexpected response from server: \n\t\t expected: %q, \n\t\t actual: %q", want, got)
	}
}

func mockStdInput(input string) (*os.File, *os.File) {
	content := []byte(input)
	tmpfile, err := ioutil.TempFile("", "test")
	if err != nil {
		log.Fatal(err)
	}

	defer os.Remove(tmpfile.Name()) // clean up

	if _, err := tmpfile.Write(content); err != nil {
		log.Fatal(err)
	}

	if _, err := tmpfile.Seek(0, 0); err != nil {
		log.Fatal(err)
	}

	originalStdin := os.Stdin
	return tmpfile, originalStdin
}
