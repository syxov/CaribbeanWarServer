package world

import (
	"CaribbeanWarServer/structs"
	"time"
)

func (self *storage) move(user *structs.User, data map[string]interface{}) {
	moveType := data["type"].(string)
	user.SetMove(moveType)
	user.Lock()
	defer user.Unlock()
	sendData := map[string]interface{}{
		"action": "move",
		"details": map[string]interface{}{
			"id":       user.ID,
			"type":     moveType,
			"location": user.Location,
			"alpha":    user.RotationAngle,
		},
	}
	for _, neigbour := range user.NearestUsers {
		neigbour.Conn.WriteJSON(sendData)
	}
	user.GetConn().WriteJSON(sendData)
}

func (self *storage) movement(user *structs.User) {
	ticker := time.NewTicker(10 * time.Millisecond)
	defer ticker.Stop()
	for user.IsInWorld() {
		timeStamp := time.Now().UnixNano()
		<-ticker.C
		self.Lock()
		beforeUpdate := user.GetPoint()
		user.UpdatePosition(float64(time.Now().UnixNano()-timeStamp) / float64(time.Second))
		afterUpdate := user.GetPoint()
		self.ocean.Update(beforeUpdate, afterUpdate)
		self.Unlock()
	}
}
