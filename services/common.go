package services

import (
	"errors"
	"github.com/gorilla/websocket"
	"sync"
)

type user struct {
	nick string
	conn *websocket.Conn
}

type WorldStruct struct {
	world map[uint]user
	sync.Mutex
}

func (self *WorldStruct) Add(id uint, nick string, conn *websocket.Conn) error {
	self.Lock()
	defer self.Unlock()
	if len(self.world) == 0 {
		self.world = make(map[uint]user)
	}
	if _, exist := self.world[id]; !exist {
		self.world[id] = user{nick: nick, conn: conn}
		return nil
	} else {
		return errors.New("User exist")
	}
}

func (self *WorldStruct) Remove(id uint) {
	self.Lock()
	defer self.Unlock()
	delete(self.world, id)
}

func (self *WorldStruct) ProcessMessage(data map[string]interface{}) {
	switch data["action"] {
	case "chat":
		self.chatMessage(data)
	}
}

func (self *WorldStruct) chatMessage(data map[string]interface{}) {
	sender := data["details"].(map[string]uint)["sender"]
	userMessage := data["details"].(map[string]string)["message"]
	message := map[string]interface{}{
		"action": "chat",
		"details": map[string]interface{}{
			"sender":  sender,
			"message": userMessage,
		},
	}
	for k, v := range self.world {
		if k != sender {
			v.conn.WriteJSON(message)
		}
	}
}
