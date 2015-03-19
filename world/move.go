package world

import (
	"CaribbeanWarServer/structs"
	"time"
)

func (self *storage) move(user *structs.User, data map[string]interface{}) {
	moveType := data["type"].(string)
	user.SetMove(moveType)
	user.Lock()
	sendData := structs.Message{"move", map[string]interface{}{
		"id":       user.ID,
		"type":     moveType,
		"location": user.Location,
		"alpha":    user.RotationAngle,
	}}
	for _, neigbour := range user.NearestUsers {
		neigbour.Conn.WriteJSON(sendData)
	}
	user.Unlock()
	user.GetConn().WriteJSON(sendData)
}

func (self *storage) movement(user *structs.User) {
	ticker := time.NewTicker(10 * time.Millisecond)
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
