/*
	layer between new users and game world
	in case successful log in push user to world
	else send a error responces
*/

package api

import (
	"github.com/gorilla/websocket"
	"net/http"
	"time"
)

type worldStr interface {
	Add(id uint, nick string, conn *websocket.Conn) error
}

var world worldStr
var db DbConnection

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func Handler(_world worldStr, _db DbConnection) func(w http.ResponseWriter, r *http.Request) {
	world = _world
	db = _db
	return func(w http.ResponseWriter, r *http.Request) {
		var data interface{}
		conn, _ := upgrader.Upgrade(w, r, nil)
		errorMessage := map[string]interface{}{"action": "fuckup"}
		if err := conn.ReadJSON(&data); err == nil {
			dataMap := data.(map[string]interface{})
			if dataMap["action"] == "auth" {
				added := auth(dataMap["details"], conn)
				if added {
					go ping(conn)
					return
				}
			} else {
				errorMessage["details"] = map[string]string{"message": "User do not logged"}
			}
		} else {
			errorMessage["details"] = map[string]string{"message": err.Error()}
		}
		conn.WriteJSON(errorMessage)
		conn.Close()
	}
}

func auth(data interface{}, conn *websocket.Conn) bool {
	added := false
	dataMap := data.(map[string]interface{})
	message := map[string]interface{}{"action": "auth"}
	if info := db.GetUserInfo(dataMap["login"].(string), dataMap["password"].(string)); info != nil {
		if err := world.Add(info.ID, info.Email, conn); err != nil {
			message["details"] = map[string]bool{
				"alreadyInGame": true,
			}
		} else {
			message["details"] = *info
			added = true
		}
	} else {
		message["details"] = "{}"
	}
	conn.WriteJSON(message)
	return added
}

func ping(conn *websocket.Conn) {
	for {
		time.Sleep(13 * time.Second)
		if err := conn.WriteMessage(websocket.TextMessage, []byte("{}")); err != nil {
			return
		}
	}
}
