package main

import (
	"github.com/gorilla/websocket"
	"net/http"
	"os"
	"runtime"
	"time"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU() * 2)
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var count byte = 0

func handler(w http.ResponseWriter, r *http.Request) {
	conn, _ := upgrader.Upgrade(w, r, nil)
	count++
	go func() {
		for {
			time.Sleep(20 * time.Second)
			if err := conn.WriteMessage(websocket.TextMessage, []byte("")); err != nil {
				return
			}
		}
	}()
	for {
		messageType, p, _ := conn.ReadMessage()
		if err := conn.WriteMessage(messageType, append(p, count)); err != nil {
			conn.Close()
			return
		}
	}
}

func main() {
	http.HandleFunc("/ws", handler)
	http.HandleFunc("/", func(writer http.ResponseWriter, r *http.Request) {
		writer.Write([]byte("Welcome me dear friend"))
	})
	http.ListenAndServe(":"+os.Getenv("PORT"), nil)
}
