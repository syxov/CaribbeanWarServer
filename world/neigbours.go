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
	user.NearestUsers = make([]structs.NearestUser, 0, len(spatials))
	for _, value := range spatials {
		converterValue := value.(*structs.User)
		user.NearestUsers = append(user.NearestUsers, structs.NearestUser{
			ID:   converterValue.ID,
			Nick: converterValue.Nick,
			Conn: converterValue.Conn,
			Ship: converterValue.SelectedShip,
		})
	}
	user.Conn.WriteJSON(map[string]interface{}{
		"action": "neigbours",
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
