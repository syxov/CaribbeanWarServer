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
	err := self.Conn.WriteJSON(message)
	self.Unlock()
	return err
}

func (self *Connection) WriteMessage(messageType int, data []byte) error {
	self.Lock()
	err := self.Conn.WriteMessage(messageType, data)
	self.Unlock()
	return err
}

func (self *Connection) ReadJSON(message interface{}) error {
	return self.Conn.ReadJSON(message)
}

func (self *Connection) Close() error {
	return self.Conn.Close()
}
