/*
	layer between new users and game world
	in case successful log in push user to world
	else send a error responces
*/

package api

import (
	"CaribbeanWarServer/harbor"
	"github.com/gorilla/websocket"
	"net/http"
	"time"
)

var db DbConnection

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func Handler(_db DbConnection) func(w http.ResponseWriter, r *http.Request) {
	db = _db
	return func(w http.ResponseWriter, r *http.Request) {
		var data interface{}
		conn, _ := upgrader.Upgrade(w, r, nil)
		go ping(conn)
		errorMessage := map[string]interface{}{"action": "fuckup"}
		if err := conn.ReadJSON(&data); err == nil {
			dataMap := data.(map[string]interface{})
			if dataMap["action"] == "auth" {
				if added := auth(dataMap["details"].(map[string]interface{}), conn); added {
					return
				} else {
					errorMessage["details"] = map[string]string{"message": "User do not added"}
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

func auth(dataMap map[string]interface{}, conn *websocket.Conn) (added bool) {
	added = false
	message := map[string]interface{}{"action": "auth"}
	if info, err := db.GetUserInfo(dataMap["login"].(string), dataMap["password"].(string)); err == nil {
		info.Conn = conn
		if err := harbor.Add(info); err == nil {
			message["details"] = map[string]interface{}{
				"authorize": true,
				"userInfo":  info,
			}
			added = true
		} else {
			message["details"] = map[string]interface{}{
				"inGame": true,
				"err":    err.Error(),
			}
		}
	} else {
		message["details"] = map[string]interface{}{
			"authorize": false,
			"details":   err.Error(),
		}
	}
	conn.WriteJSON(message)
	return
}

func ping(conn *websocket.Conn) {
	for {
		time.Sleep(13 * time.Second)
		if err := conn.WriteMessage(websocket.TextMessage, []byte("{}")); err != nil {
			return
		}
	}
}
