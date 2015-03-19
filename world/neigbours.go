package world

import (
	"CaribbeanWarServer/structs"
	"time"
)

const radius = 100000

func (self *storage) findNeigbours(user *structs.User) {
	user.Lock()
	if user.NearestUsers == nil {
		user.NearestUsers = make([]structs.NearestUser, 0, 5)
	}
	rect := user.Bounds(radius)
	user.Unlock()
	self.Lock()
	spatials := self.ocean.Search(rect)
	self.Unlock()
	user.Lock()
	defer user.Unlock()
	nearestUsers := make([]structs.NearestUser, 0, len(spatials))
	for _, value := range spatials {
		convertedValue := value.Data().(*structs.User)
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
	addedGamersChanel, removedGamersChanel := make(chan []structs.NearestUser), make(chan []structs.NearestUser)
	go getDifference(&nearestUsers, &user.NearestUsers, addedGamersChanel)
	go getDifference(&user.NearestUsers, &nearestUsers, removedGamersChanel)
	addedGamers, removedGamers := <-addedGamersChanel, <-removedGamersChanel
	if len(addedGamers) != 0 || len(removedGamers) != 0 {
		user.NearestUsers = nearestUsers
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

func getDifference(p_firstSlice, p_secondSlice *[]structs.NearestUser, channel chan []structs.NearestUser) {
	firstSlice := *p_firstSlice
	secondSlice := *p_secondSlice
	difference := make([]structs.NearestUser, 0, 10)
	for _, firstSliceUser := range firstSlice {
		isShouldBeAddToDiff := true
		for _, secondSliceUser := range secondSlice {
			if firstSliceUser.ID == secondSliceUser.ID {
				isShouldBeAddToDiff = false
				break
			}
		}
		if isShouldBeAddToDiff {
			difference = append(difference, firstSliceUser)
		}
	}
	channel <- difference
}
