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
		var data map[string]interface{}
		conn, _ := upgrader.Upgrade(w, r, nil)
		errorMessage := map[string]interface{}{"action": "fuckup"}
		if err := conn.ReadJSON(&data); err == nil {
			if data["action"] == "auth" {
				if added := auth(data["details"].(map[string]string), conn); added {
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

func auth(data map[string]string, conn *websocket.Conn) bool {
	added := false
	message := map[string]interface{}{"action": "auth"}
	if info := db.GetUserInfo(data["login"], data["password"]); info != nil {
		if err := world.Add(info["id"].(uint), info["nick"].(string), conn); err != nil {
			message["details"] = map[string]bool{
				"inGame": true,
			}
		} else {
			message["details"] = map[string]interface{}{"authorize": true}
			for key, value := range info {
				message["details"].(map[string]interface{})[key] = value
			}
			added = true
		}
	} else {
		message["details"] = map[string]bool{"authorize": false}
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
