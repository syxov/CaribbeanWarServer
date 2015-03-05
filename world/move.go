package world

import (
	"CaribbeanWarServer/structs"
	"time"
)

func (self *storage) move(user *structs.User, data map[string]interface{}) {
	moveType := data["type"].(string)
	sendData := map[string]interface{}{
		"action": "move",
		"details": map[string]interface{}{
			"id":       user.ID,
			"type":     moveType,
			"location": user.Location,
			"alpha":    user.RotationAngle,
		},
	}
	user.SetMove(moveType)
	user.Lock()
	for _, neigbour := range user.NearestUsers {
		neigbour.Conn.WriteJSON(sendData)
	}
	user.GetConn().WriteJSON(sendData)
	user.Unlock()
}

func (self *storage) movement(user *structs.User) {
	ticker := time.NewTicker(10 * time.Millisecond)
	defer ticker.Stop()
	for user.IsInWorld() {
		<-ticker.C
		self.Lock()
		isDeleted := self.ocean.Delete(user)
		if isDeleted {
			user.UpdatePosition(0.01)
			self.ocean.Insert(user)
		}
		self.Unlock()
	}
}
