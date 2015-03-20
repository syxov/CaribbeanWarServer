package world

import (
	"CaribbeanWarServer/api"
	"CaribbeanWarServer/rtree"
	"CaribbeanWarServer/structs"
)

type storage struct {
	ocean *rtree.Rtree
	db    api.DbConnection
}

var world storage
var addToHarbor func(*structs.User) error

func init() {
	world.ocean = rtree.NewTree(2, 2, 50)
	world.db = api.DbConnection{}
	world.db.Open()
}

func InitHarbor(add func(*structs.User) error) {
	addToHarbor = add
}

func Add(user *structs.User) {
	world.add(user)
}

func (self *storage) add(user *structs.User) {
	user.SetIsInWorld(true)
	self.ocean.Insert(user)
	self.findNeigbours(user)
	go self.message(user)
	go self.findNeigboursRepeater(user)
	go self.movement(user)

}

func (self *storage) remove(user *structs.User, needAddToHarbor bool) {
	user.Lock()
	user.NearestUsers = nil
	user.SelectedShip = nil
	user.SetIsInWorld(false)
	self.ocean.Delete(user)
	self.db.SaveUserLocation(user)
	if needAddToHarbor {
		addToHarbor(user)
	}
	user.Unlock()
}
