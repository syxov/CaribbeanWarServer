package api

import (
	"github.com/gorilla/websocket"
	"net/http"
	"time"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func Handler(w http.ResponseWriter, r *http.Request) {
	conn, _ := upgrader.Upgrade(w, r, nil)
	defer conn.Close()
	go ping(conn)
	for {
		var data interface{}
		if err := conn.ReadJSON(&data); err == nil {
			dataMap := data.(map[string]interface{})
			switch dataMap["action"] {
			case "auth":
				auth(dataMap["details"], conn)
			}
		} else {
			errorMessage := map[string]interface{}{"action": "fuckup"}
			errorMessage["details"] = map[string]string{"message": err.Error()}
			send(conn, errorMessage)
		}
	}
}

func auth(data interface{}, conn *websocket.Conn) {
	dataMap := data.(map[string]interface{})
	message := map[string]interface{}{"action": "auth"}
	if DbConn.CheckUserExist(dataMap["email"].(string), dataMap["password"].(string)) {
		message["details"] = map[string]string{"result": "you are awesome"}
	} else {
		message["details"] = map[string]string{"result": "fuck you"}
	}
	send(conn, message)
}

func ping(conn *websocket.Conn) {
	for {
		time.Sleep(10 * time.Second)
		if err := conn.WriteMessage(websocket.TextMessage, []byte{}); err != nil {
			panic("Cannot send message")
		}
	}
}

func send(conn *websocket.Conn, data interface{}) {
	if err := conn.WriteJSON(data); err != nil {
		panic("Cannot send message")
	}
}
