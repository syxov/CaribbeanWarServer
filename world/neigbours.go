package world

import (
	"CaribbeanWarServer/commonStructs"
	"CaribbeanWarServer/messagesStructs"
	"CaribbeanWarServer/structs"
	"time"
)

const radius = 100000

func (self *storage) findNeigbours(user *structs.User) {
	user.Lock()
	if user.NearestUsers == nil {
		user.NearestUsers = make([]commonStructs.NearestUser, 0, 5)
	}
	rect := user.Bounds(radius)
	user.Unlock()
	spatials := self.ocean.SearchIntersect(rect)
	nearestUsers := make([]commonStructs.NearestUser, 0, len(spatials))
	user.Lock()
	for _, value := range spatials {
		convertedValue := value.(*structs.User)
		if convertedValue.ID != user.ID {
			nearestUsers = append(nearestUsers, commonStructs.NearestUser{
				ID:            convertedValue.ID,
				Conn:          convertedValue.GetConn(),
				Ship:          convertedValue.SelectedShip,
				Nick:          convertedValue.Nick,
				Location:      convertedValue.Location,
				RotationAngle: convertedValue.RotationAngle,
			})
		}
	}
	addedGamersChanel, removedGamersChanel := make(chan []commonStructs.NearestUser), make(chan []commonStructs.NearestUser)
	go getDifference(&nearestUsers, &user.NearestUsers, addedGamersChanel)
	go getDifference(&user.NearestUsers, &nearestUsers, removedGamersChanel)
	addedGamers, removedGamers := <-addedGamersChanel, <-removedGamersChanel
	if len(addedGamers) != 0 || len(removedGamers) != 0 {
		user.NearestUsers = nearestUsers
		user.GetConn().WriteJSON(messagesStructs.Message{"neighbours", map[string]interface{}{
			"added":   addedGamers,
			"removed": removedGamers,
		}})
	}
	user.Unlock()
}

func (self *storage) findNeigboursRepeater(user *structs.User) {
	for user.IsInWorld() {
		self.findNeigbours(user)
		time.Sleep(time.Second)
	}
}

func getDifference(p_firstSlice, p_secondSlice *[]commonStructs.NearestUser, channel chan []commonStructs.NearestUser) {
	firstSlice := *p_firstSlice
	secondSlice := *p_secondSlice
	difference := make([]commonStructs.NearestUser, 0, 10)
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
