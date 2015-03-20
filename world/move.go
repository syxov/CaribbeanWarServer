package world

import (
	"CaribbeanWarServer/messagesStructs"
	"CaribbeanWarServer/structs"
	"time"
)

func (self *storage) move(user *structs.User, data messagesStructs.MoveIncome) {
	user.SetMove(data.Details.Type)
	user.Lock()
	sendData := messagesStructs.MoveOutcome{
		Action: "move",
		Details: messagesStructs.MoveOutcomeDetails{
			ID:       user.ID,
			Type:     data.Details.Type,
			Location: user.Location,
			Alpha:    user.RotationAngle,
		},
	}
	for _, neigbour := range user.NearestUsers {
		neigbour.Conn.WriteJSON(sendData)
	}
	user.GetConn().WriteJSON(sendData)
	user.Unlock()
}

func (self *storage) movement(user *structs.User) {
	ticker := time.NewTicker(1000 / 60 * time.Millisecond)
	for user.IsInWorld() {
		timeStamp := time.Now().UnixNano()
		<-ticker.C
		isDeleted := self.ocean.Delete(user)
		if isDeleted {
			user.UpdatePosition(float64(time.Now().UnixNano()-timeStamp) / float64(time.Second))
			self.ocean.Insert(user)
		}
	}
	ticker.Stop()
}
