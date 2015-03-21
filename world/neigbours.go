package world

import (
	"CaribbeanWarServer/commonStructs"
	"CaribbeanWarServer/messagesStructs"
	"CaribbeanWarServer/structs"
	"sort"
	"time"
)

const radius = 100000

func (self *storage) findNeigbours(user *structs.User) {
	user.Lock()
	if user.NearestUsers == nil {
		user.NearestUsers = make(commonStructs.NearestUsers, 0)
	}
	rect := user.Bounds(radius)
	user.Unlock()
	spatials := self.ocean.SearchIntersect(rect)
	nearestUsers := make(commonStructs.NearestUsers, len(spatials))
	for index, value := range spatials {
		convertedValue := value.(*structs.User)
		if convertedValue.ID != user.ID {
			nearestUsers[index] = commonStructs.NearestUser{
				ID:            convertedValue.ID,
				Conn:          convertedValue.GetConn(),
				Ship:          convertedValue.SelectedShip,
				Nick:          convertedValue.Nick,
				Location:      convertedValue.Location,
				RotationAngle: convertedValue.RotationAngle,
			}
		}
	}
	sort.Sort(&nearestUsers)
	addedGamersChanel, removedGamersChanel := make(chan commonStructs.NearestUsers), make(chan commonStructs.NearestUsers)
	go getDifference(&nearestUsers, &user.NearestUsers, addedGamersChanel)
	go getDifference(&user.NearestUsers, &nearestUsers, removedGamersChanel)
	addedGamers, removedGamers := <-addedGamersChanel, <-removedGamersChanel
	if len(addedGamers) != 0 || len(removedGamers) != 0 {
		user.Lock()
		user.NearestUsers = nearestUsers
		user.GetConn().WriteJSON(messagesStructs.Message{"neighbours", map[string]interface{}{
			"added":   addedGamers,
			"removed": removedGamers,
		}})
		user.Unlock()
	}
	close(addedGamersChanel)
	close(removedGamersChanel)
}

func (self *storage) findNeigboursRepeater(user *structs.User) {
	for user.IsInWorld() {
		self.findNeigbours(user)
		time.Sleep(2 * time.Second)
	}
}

func getDifference(p_firstSlice, p_secondSlice *commonStructs.NearestUsers, channel chan commonStructs.NearestUsers) {
	firstSlice := *p_firstSlice
	secondSlice := *p_secondSlice
	difference := make(commonStructs.NearestUsers, 0, 3)
	for _, firstSliceUser := range firstSlice {
		if index := sort.Search(len(secondSlice), func(j int) bool { return firstSliceUser.ID == secondSlice[j].ID }); index >= 0 && index < len(secondSlice) {
			difference = append(difference, firstSliceUser)
		}
	}
	channel <- difference
}
