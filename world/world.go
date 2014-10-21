package world

import (
	"CaribbeanWarServer/structs"
	"errors"
	"github.com/dhconnelly/rtreego"
	"sync"
)

type storage struct {
	ocean    *rtreego.Rtree
	userList map[uint]*structs.User
	sync.Mutex
}

var world storage

func Add(user *structs.User) {
	world.add(user)
}

func (self *storage) add(user *structs.User) {
	self.init()
	self.userList[user.ID] = user
	self.ocean.Insert(&node{ID: user.ID, Location: &user.Location})
	go self.message(user)
}

func (self *storage) remove(user *structs.User) {
	self.Lock()
	defer self.Unlock()
	delete(self.userList, user.ID)
	self.ocean.Delete(&node{ID: user.ID, Location: &user.Location})
}

func (self *storage) init() {
	self.Lock()
	defer self.Unlock()
	if self.ocean == nil {
		self.ocean = rtreego.NewTree(2, 2, 10)
		self.userList = make(map[uint]*structs.User, 1000)
	}
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
			details := json["Details"].(map[string]interface{})
			switch json["Action"] {
			case "chat":
				self.chat(details)
			}
		} else {
			if err.Error() == "EOF" {
				self.remove(user)
				return
			} else {
				user.Conn.WriteJSON(map[string]string{
					"Action":  "fuckup",
					"Details": err.Error(),
				})
			}
		}
	}
}
