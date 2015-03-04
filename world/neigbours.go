package world

import (
	"CaribbeanWarServer/rtree"
	"CaribbeanWarServer/structs"
	"time"
)

const radius = 100000

func (self *storage) findNeigbours(user *structs.User) {
	user.Lock()
	self.Lock()
	defer self.Unlock()
	defer user.Unlock()
	rect, _ := rtree.NewRect(rtree.Point{user.Location.X, user.Location.Y}, []float64{radius, radius})
	if user.NearestUsers == nil {
		user.NearestUsers = make([]structs.NearestUser, 0, 5)
	}
	spatials := self.ocean.SearchIntersect(rect)
	nearestUsers := make([]structs.NearestUser, 0, len(spatials))
	for _, value := range spatials {
		convertedValue := value.(*structs.User)
		if convertedValue.ID != user.ID {
			nearestUsers = append(nearestUsers, structs.NearestUser{
				ID:            &convertedValue.ID,
				Conn:          convertedValue.GetConn(),
				Ship:          convertedValue.SelectedShip,
				Nick:          &convertedValue.Nick,
				Location:      convertedValue.Location,
				RotationAngle: convertedValue.RotationAngle,
			})
		}
	}
	addedGamers := getAddedGamers(&user.NearestUsers, &nearestUsers)
	removedGamers := getRemovedGamers(&user.NearestUsers, &nearestUsers)
	user.NearestUsers = nearestUsers
	if len(addedGamers) != 0 || len(removedGamers) != 0 {
		user.GetConn().WriteJSON(map[string]interface{}{
			"action": "neighbours",
			"details": map[string][]structs.NearestUser{
				"added":   addedGamers,
				"removed": removedGamers,
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

func getAddedGamers(p_oldNearestUsers, p_newNearestUsers *[]structs.NearestUser) []structs.NearestUser {
	oldNearestUsers := *p_oldNearestUsers
	newNearestUsers := *p_newNearestUsers
	newGamersSlice := make([]structs.NearestUser, 0, 2)
	for _, nearestUser := range newNearestUsers {
		isNew := true
		for _, oldNearestUser := range oldNearestUsers {
			if nearestUser.ID == oldNearestUser.ID {
				isNew = false
				break
			}
		}
		if isNew {
			newGamersSlice = append(newGamersSlice, nearestUser)
		}
	}
	return newGamersSlice
}

func getRemovedGamers(p_oldNearestUsers, p_newNearestUsers *[]structs.NearestUser) []structs.NearestUser {
	oldNearestUsers := *p_oldNearestUsers
	newNearestUsers := *p_newNearestUsers
	removedGamersSlice := make([]structs.NearestUser, 0, 2)
	for _, oldNearestUser := range oldNearestUsers {
		isShouldBeRemoved := true
		for _, newNearestUser := range newNearestUsers {
			if oldNearestUser.ID == newNearestUser.ID {
				isShouldBeRemoved = false
				break
			}
		}
		if isShouldBeRemoved {
			removedGamersSlice = append(removedGamersSlice, oldNearestUser)
		}
	}
	return removedGamersSlice
}
