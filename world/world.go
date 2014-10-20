package world

import (
	"CaribbeanWarServer/structs"
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
	self.ocean.Insert(&node{ID: user.ID, Conn: user.Conn, Location: &user.Location})
}

func (self *storage) init() {
	self.Lock()
	defer self.Unlock()
	if self.ocean == nil {
		self.ocean = rtreego.NewTree(2, 2, 10)
		self.userList = make(map[uint]*structs.User, 1000)
	}
}
