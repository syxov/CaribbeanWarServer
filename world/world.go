package world

import (
	"CaribbeanWarServer/api"
	"CaribbeanWarServer/quadtree"
	"CaribbeanWarServer/structs"
	"sync"
)

type storage struct {
	ocean *quadtree.QuadTree
	db    api.DbConnection
	sync.Mutex
}

var world storage
var addToHarbor func(*structs.User) error

func init() {
	bound := quadtree.NewAABB(&quadtree.Point{X: 0, Y: 0}, &quadtree.Point{X: 10000000, Y: 100000000})
	world.ocean = quadtree.New(bound, 0, nil)
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
	self.Lock()
	self.ocean.Insert(user.GetPoint())
	self.Unlock()
	self.findNeigbours(user)
	go self.message(user)
	go self.findNeigboursRepeater(user)
	go self.movement(user)

}

func (self *storage) remove(user *structs.User, needAddToHarbor bool) {
	self.Lock()
	user.Lock()
	defer user.Unlock()
	defer self.Unlock()
	user.NearestUsers = nil
	user.SelectedShip = nil
	user.SetIsInWorld(false)
	self.ocean.Remove(user.GetPoint())
	self.db.SaveUserLocation(user)
	if needAddToHarbor {
		addToHarbor(user)
	}
}
