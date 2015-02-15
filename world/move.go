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
	user.Lock()
	user.SetMove(moveType)
	for _, neigbour := range user.NearestUsers {
		neigbour.Conn.WriteJSON(sendData)
	}
	user.Unlock()
	user.GetConn().WriteJSON(sendData)
}

func (self *storage) movement(user *structs.User) {
	ticker := time.NewTicker(10 * time.Millisecond)
	defer ticker.Stop()
	for user.IsInWorld() {
		tick := <-ticker.C
		self.Lock()
		isDeleted := self.ocean.Delete(user)
		if isDeleted {
			user.Lock()
			user.UpdatePosition(float64(tick.Nanosecond() / int(time.Second)))
			user.Unlock()
			self.ocean.Insert(user)
		}
		self.Unlock()
	}
}
