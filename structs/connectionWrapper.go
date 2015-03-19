package structs

import (
	"github.com/gorilla/websocket"
	"github.com/pquerna/ffjson/ffjson"
	"sync"
)

type Connection struct {
	Conn *websocket.Conn
	sync.Mutex
}

func (self *Connection) WriteJSON(message interface{}) error {
	self.Lock()
	defer self.Unlock()
	if json, err := ffjson.Marshal(message); err == nil {
		return self.Conn.WriteMessage(websocket.TextMessage, json)
	} else {
		return err
	}
}

func (self *Connection) WriteMessage(messageType int, data []byte) error {
	self.Lock()
	defer self.Unlock()
	return self.Conn.WriteMessage(messageType, data)
}

func (self *Connection) ReadJSON(message interface{}) error {
	return self.Conn.ReadJSON(message)
}

func (self *Connection) Close() error {
	self.Lock()
	defer self.Unlock()
	return self.Conn.Close()
}
