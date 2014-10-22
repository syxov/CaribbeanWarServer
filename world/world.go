package world

import (
	"CaribbeanWarServer/rtree"
	"CaribbeanWarServer/structs"
	"errors"
	"sync"
)

type storage struct {
	ocean *rtree.Rtree
	sync.Mutex
}

var world storage

func init() {
	world.ocean = rtree.NewTree(2, 2, 10)
}

func Add(user *structs.User) {
	world.add(user)
}

func (self *storage) add(user *structs.User) {
	self.Lock()
	defer self.Unlock()
	self.ocean.Insert(user)
	go self.message(user)
}

func (self *storage) remove(user *structs.User) {
	self.Lock()
	defer self.Unlock()
	self.ocean.Delete(user)
}

func (self *storage) message(user *structs.User) {
	defer func() {
		if err := recover(); err != nil {
			var message string
			switch err.(type) {
			case error:
				message = err.(error).Error()
			default:
				message = "Something wrong"
			}
			user.Conn.WriteJSON(errors.New(message))
			user.Conn.Close()
			self.remove(user)
		}
	}()

	var message interface{}
	for {
		if err := user.Conn.ReadJSON(&message); err == nil {
			json := message.(map[string]interface{})
			details := json["details"].(map[string]interface{})
			switch json["action"] {
			case "chat":
				self.chat(details)
			}
		} else {
			if err.Error() == "EOF" {
				self.remove(user)
				return
			} else {
				user.Conn.WriteJSON(map[string]string{
					"action":  "fuckup",
					"details": err.Error(),
				})
			}
		}
	}
}
