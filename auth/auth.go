/*
	layer between new users and game world
	in case of successful login push user to world
	else send a error responce
*/

package auth

import (
	"CaribbeanWarServer/messagesStructs"
	"CaribbeanWarServer/structs"
	"github.com/gorilla/websocket"
	"net/http"
	"time"
)

type DbConnection interface {
	GetUserInfo(string, string) (*structs.User, error)
}

type Harbor interface {
	Add(*structs.User) error
}

var db DbConnection
var harbor Harbor

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func Handler(_db DbConnection, _harbor Harbor) func(w http.ResponseWriter, r *http.Request) {
	db = _db
	harbor = _harbor
	return func(w http.ResponseWriter, r *http.Request) {
		_conn, _ := upgrader.Upgrade(w, r, nil)
		conn := &structs.Connection{Conn: _conn}
		defer func() {
			if err := recover(); err != nil {
				conn.WriteJSON(err)
				conn.Close()
			}
		}()
		var data messagesStructs.Message
		go ping(conn)
		if err := conn.ReadJSON(&data); err != nil {
			panic(messagesStructs.ErrorMessage(err.Error()))
		}
		if data.Action != "auth" {
			panic(messagesStructs.ErrorMessage("User do not logged"))
		}
		if authorized := auth(data.Details.(map[string]interface{}), conn); !authorized {
			panic(messagesStructs.ErrorMessage("User do not added"))
		}
	}
}

func auth(dataMap map[string]interface{}, conn *structs.Connection) bool {
	added := false
	message := map[string]interface{}{"action": "auth"}
	if info, err := db.GetUserInfo(dataMap["login"].(string), dataMap["password"].(string)); err == nil {
		info.SetConn(conn)
		if err := harbor.Add(info); err == nil {
			message["details"] = map[string]interface{}{
				"authorize": true,
				"userInfo":  info,
			}
			added = true
		} else {
			message["details"] = map[string]interface{}{
				"authorize": false,
				"inGame":    true,
			}
		}
	} else {
		message["details"] = map[string]interface{}{
			"authorize": false,
			"details":   err.Error(),
		}
	}
	conn.WriteJSON(message)
	return added
}

var message []byte = []byte("{}")

func ping(conn *structs.Connection) {
	for conn.WriteMessage(websocket.TextMessage, message) == nil {
		time.Sleep(20 * time.Second)
	}
}
