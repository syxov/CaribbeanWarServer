// CaribbeanWar project main.go
package main

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"os"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func handler(w http.ResponseWriter, r *http.Request) {
	if conn, err := upgrader.Upgrade(w, r, nil); err == nil {
		for {
			messageType, p, err := conn.ReadMessage()
			if err == nil {
				conn.WriteMessage(messageType, p)
			} else {
				conn.WriteMessage(messageType, []byte("Fuck"))
			}
		}
	} else {
		log.Print(err)
	}
}

func main() {
	log.Print("Server started")
	http.HandleFunc("/ws", handler)
	http.HandleFunc("/", func(writer http.ResponseWriter, r *http.Request) {
		writer.Write([]byte("some"))
	})
	http.ListenAndServe(":"+os.Getenv("PORT"), nil)
}
