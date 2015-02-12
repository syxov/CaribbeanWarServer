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
	nearestUsers := make([]structs.NearestUser, 0, len(spatials))
	for _, value := range spatials {
		converterValue := value.(*structs.User)
		if converterValue.ID != user.ID {
			nearestUsers = append(nearestUsers, structs.NearestUser{
				ID:   converterValue.ID,
				Conn: converterValue.GetConn(),
				Ship: converterValue.SelectedShip,
			})
		}
	}
	if listChanged(nearestUsers, user.NearestUsers) {
		user.NearestUsers = nearestUsers
		user.GetConn().WriteJSON(map[string]interface{}{
			"action": "neigbours",
			"details": map[string]interface{}{
				"users": user.NearestUsers,
			},
		})
	}
}

func (self *storage) findNeigboursRepeater(user *structs.User) {
	for user.IsInWorld() {
		self.findNeigbours(user)
		time.Sleep(time.Second)
	}
}

func listChanged(a, b []structs.NearestUser) bool {
	if len(a) != len(b) {
		return true
	}
	lenB := len(b)
	for _, value := range a {
		for index, innerValue := range b {
			if innerValue.ID == value.ID {
				break
			} else if index == lenB {
				return true
			}
		}
	}
	return false
}
