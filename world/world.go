package world

import (
	"errors"
	"github.com/gorilla/websocket"
	"sync"
)

type user struct {
	nick string
	conn *websocket.Conn
}

type DbConnection interface {
	GetShopItems(id uint)
}

type WorldStruct struct {
	world  map[uint]user
	DbConn DbConnection
	sync.Mutex
}

func (self *WorldStruct) Add(id uint, nick string, conn *websocket.Conn) error {
	self.Lock()
	defer self.Unlock()
	if len(self.world) == 0 {
		self.world = make(map[uint]user)
	}
	if _, exist := self.world[id]; exist {
		return errors.New("User exist")
	}
	self.world[id] = user{nick: nick, conn: conn}
	go self.processMessage(id, conn)
	return nil
}

func (self *WorldStruct) Remove(id uint) {
	self.Lock()
	defer self.Unlock()
	delete(self.world, id)
}

func (self *WorldStruct) processMessage(id uint, conn *websocket.Conn) {
	defer func() {
		if err := recover(); err != nil {
			self.sendErrorMessage(conn, err)
			conn.Close()
			self.Remove(id)
		}
	}()
	for {
		var data interface{}
		if err := conn.ReadJSON(&data); err == nil {
			convertedData := data.(map[string]interface{})
			switch convertedData["action"] {
			case "chat":
				self.chat(convertedData)
			case "shop":
				self.shop(convertedData)
			}
		} else {
			if err.Error() == "EOF" { //if connection closed
				self.Remove(id)
				return
			} else { //Problem with json converting
				self.sendErrorMessage(conn, err.Error())
			}
		}
	}
}

func (self *WorldStruct) sendErrorMessage(conn *websocket.Conn, err interface{}) {
	errorMessage := map[string]interface{}{"action": "fuckup"}
	errorMessage["details"] = map[string]interface{}{"message": err}
	conn.WriteJSON(errorMessage)
}
