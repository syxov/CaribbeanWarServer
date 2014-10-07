package api

import (
	"CaribbeanWarServer/services"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

var World services.WorldStruct

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func Handler(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			log.Print(err)
		}
	}()
	conn, _ := upgrader.Upgrade(w, r, nil)
	defer conn.Close()
	go ping(conn)
	for {
		var data interface{}
		if err := conn.ReadJSON(&data); err == nil {
			dataMap := data.(map[string]interface{})
			if dataMap["action"] == "auth" {
				if id := auth(dataMap["details"], conn); id != 0 {
					defer func() {
						World.Remove(id)
					}()
				} else {
					return
				}
			} else {
				World.ProcessMessage(dataMap)
			}
		} else {
			errorMessage := map[string]interface{}{"action": "fuckup"}
			errorMessage["details"] = map[string]string{"message": err.Error()}
			send(conn, errorMessage)
		}
	}
}

func auth(data interface{}, conn *websocket.Conn) uint {
	dataMap := data.(map[string]interface{})
	message := map[string]interface{}{"action": "auth"}
	returnValue := uint(0)
	if info := DbConn.GetUserInfo(dataMap["email"].(string), dataMap["password"].(string)); info != nil {
		if err := World.Add(info.ID, info.Email, conn); err != nil {
			message["details"] = map[string]bool{
				"alreadyInGame": true,
			}
		} else {
			message["details"] = *info
			returnValue = info.ID
		}
	} else {
		message["details"] = "User not found"
	}
	send(conn, message)
	return returnValue
}

func ping(conn *websocket.Conn) {
	for {
		time.Sleep(15 * time.Second)
		if err := conn.WriteMessage(websocket.TextMessage, []byte("{}")); err != nil {
			return
		}
	}
}

func send(conn *websocket.Conn, data interface{}) {
	if err := conn.WriteJSON(data); err != nil {
		panic(err)
	}
}
