package world

import (
	"CaribbeanWarServer/rtree"
	"CaribbeanWarServer/structs"
	"sync"
)

type storage struct {
	ocean *rtree.Rtree
	sync.Mutex
}

var world storage
var addToHarbor func(*structs.User) error

func init() {
	world.ocean = rtree.NewTree(2, 2, 50)
}

func InitHarbor(add func(*structs.User) error) {
	addToHarbor = add
}

func Add(user *structs.User) {
	world.add(user)
}

func (self *storage) add(user *structs.User) {
	user.InWorld = true
	self.Lock()
	self.ocean.Insert(user)
	self.Unlock()
	self.findNeigbours(user)
	go self.message(user)
	go self.findNeigboursRepeater(user)
	go self.movement(user)
}

func (self *storage) remove(user *structs.User) {
	self.Lock()
	defer self.Unlock()
	user.NearestUsers = nil
	user.SelectedShip = nil
	user.InWorld = false
	self.ocean.Delete(user)
	addToHarbor(user)
}
