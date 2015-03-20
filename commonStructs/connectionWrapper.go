package commonStructs

import (
	"github.com/gorilla/websocket"
	"sync"
)

type Connection struct {
	Conn *websocket.Conn
	sync.Mutex
}

func (self *Connection) WriteJSON(message interface{}) error {
	self.Lock()
	defer self.Unlock()
	return self.Conn.WriteJSON(message)
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
	return self.Conn.Close()
}
