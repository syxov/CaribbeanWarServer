/*
	layer between new users and game world
	in case successful log in push user to world
	else send a error responces
*/

package api

import (
	"CaribbeanWarServer/structs"
	"github.com/gorilla/websocket"
	"net/http"
)

type harborStr interface {
	Add(data *structs.User) error
}

var harbor harborStr
var db DbConnection

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func Handler(_harbor harborStr, _db DbConnection) func(w http.ResponseWriter, r *http.Request) {
	harbor = _harbor
	db = _db
	return func(w http.ResponseWriter, r *http.Request) {
		var data interface{}
		conn, _ := upgrader.Upgrade(w, r, nil)
		errorMessage := map[string]interface{}{"Action": "fuckup"}
		if err := conn.ReadJSON(&data); err == nil {
			dataMap := data.(map[string]interface{})
			if dataMap["Action"] == "auth" {
				if added := auth(dataMap["Details"].(map[string]interface{}), conn); added {
					return
				} else {
					errorMessage["Details"] = map[string]string{"Message": "User do not added"}
				}
			} else {
				errorMessage["Details"] = map[string]string{"Message": "User do not logged"}
			}
		} else {
			errorMessage["Details"] = map[string]string{"Message": err.Error()}
		}
		conn.WriteJSON(errorMessage)
		conn.Close()
	}
}

func auth(dataMap map[string]interface{}, conn *websocket.Conn) (added bool) {
	added = false
	message := map[string]interface{}{"Action": "auth"}
	if info, err := db.GetUserInfo(dataMap["Login"].(string), dataMap["Password"].(string)); err == nil {
		info.Conn = conn
		if err := harbor.Add(info); err == nil {
			message["Details"] = map[string]interface{}{
				"Authorize": true,
				"UserInfo":  info,
			}
			added = true
		} else {
			message["Details"] = map[string]interface{}{
				"InGame":    true,
				"Error":     err.Error(),
				"Authorize": false,
			}
		}
	} else {
		message["Details"] = map[string]interface{}{
			"Authorize": false,
			"Details":   err.Error(),
		}
	}
	conn.WriteJSON(message)
	return
}
