package world

import (
	"CaribbeanWarServer/rtree"
	"CaribbeanWarServer/structs"
	"time"
)

func (self *storage) findNeigbours(user *structs.User) {
	rect, _ := rtree.NewRect(rtree.Point{user.Location.X, user.Location.Y}, []float64{100, 100})
	self.Lock()
	defer self.Unlock()
	spatials := self.ocean.SearchIntersect(rect)
	user.NearestUsers = make([]*structs.User, 0, len(spatials))
	for _, value := range spatials {
		user.NearestUsers = append(user.NearestUsers, value.(*structs.User))
	}
}

func (self *storage) findNeigboursRepeater(user *structs.User) {
	for {
		if userGoToHarbor := user.SelectedShip == nil; userGoToHarbor {
			return
		}
		self.findNeigbours(user)
		time.Sleep(3 * 1000)
	}
}
