package world

import (
	"CaribbeanWarServer/rtree"
	"CaribbeanWarServer/structs"
	"time"
)

const radius = 100

func (self *storage) findNeigbours(user *structs.User) {
	rect, _ := rtree.NewRect(rtree.Point{user.Location.X, user.Location.Y}, []float64{radius, radius})
	self.Lock()
	defer self.Unlock()
	spatials := self.ocean.SearchIntersect(rect)
	user.NearestUsers = make([]*structs.User, 0, len(spatials))
	for _, value := range spatials {
		user.NearestUsers = append(user.NearestUsers, value.(*structs.User))
	}
	user.Conn.WriteJSON(map[string]interface{}{
		"action": "nearestUsers",
		"details": map[string]interface{}{
			"users": user.NearestUsers,
		},
	})
}

func (self *storage) findNeigboursRepeater(user *structs.User) {
	for user.InWorld {
		self.findNeigbours(user)
		time.Sleep(time.Second)
	}
}
