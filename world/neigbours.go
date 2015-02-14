package world

import (
	"CaribbeanWarServer/rtree"
	"CaribbeanWarServer/structs"
	"sync"
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
		convertedValue := value.(*structs.User)
		if convertedValue.ID != user.ID {
			nearestUsers = append(nearestUsers, structs.NearestUser{
				ID:   &convertedValue.ID,
				Conn: convertedValue.GetConn(),
				Ship: convertedValue.SelectedShip,
				Nick: &convertedValue.Nick,
			})
		}
	}
	group := &sync.WaitGroup{}
	group.Add(2)
	var (
		addedGamers, removedGamers *[]structs.NearestUser
	)
	go getAddedGamers(&user.NearestUsers, &nearestUsers, addedGamers, group)
	go getRemovedGamers(&user.NearestUsers, &nearestUsers, removedGamers, group)
	group.Wait()
	if len(*addedGamers) != 0 || len(*removedGamers) != 0 {
		user.GetConn().WriteJSON(map[string]interface{}{
			"action": "nieghbours",
			"details": map[string]*[]structs.NearestUser{
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

func getAddedGamers(p_oldNearestUsers, p_newNearestUsers, newGamers *[]structs.NearestUser, waitGroup *sync.WaitGroup) {
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
	newGamers = &newGamersSlice
	waitGroup.Done()
}

func getRemovedGamers(p_oldNearestUsers, p_newNearestUsers, removedGamers *[]structs.NearestUser, waitGroup *sync.WaitGroup) {
	oldNearestUsers := *p_oldNearestUsers
	newNearestUsers := *p_newNearestUsers
	removedGamersSlice := make([]structs.NearestUser, 0, 2)
	for _, nearestUser := range oldNearestUsers {
		isShouldBeRemoved := true
		for _, oldNearestUser := range newNearestUsers {
			if nearestUser.ID == oldNearestUser.ID {
				isShouldBeRemoved = false
				break
			}
		}
		if isShouldBeRemoved {
			removedGamersSlice = append(removedGamersSlice, nearestUser)
		}
	}
	removedGamers = &removedGamersSlice
	waitGroup.Done()
}
