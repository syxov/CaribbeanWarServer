package world

import (
	"CaribbeanWarServer/messagesStructs"
	"CaribbeanWarServer/structs"
	"time"
)

func (self *storage) move(user *structs.User, ch chan *messagesStructs.MoveIncome) {
	for {
		if data, ok := <-ch; ok {
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
		} else {
			return
		}
	}
}

func (self *storage) movement(user *structs.User) {
	ticker := time.NewTicker(1000 / 35 * time.Millisecond)
	for user.IsInWorld() {
		if user.IsKilled() {
			time.Sleep(2 * time.Second)
			continue
		}
		timeStamp := time.Now().UnixNano()
		<-ticker.C
		if speed := user.UpdateSpeed(float64(time.Now().UnixNano()-timeStamp) / float64(time.Second)); speed > 0.00001 {
			isDeleted := self.ocean.Delete(user)
			if isDeleted {
				user.UpdatePosition()
				self.ocean.Insert(user)
			}
		}
		user.GetConn().WriteJSON(messagesStructs.PositionMessage{
			Message: messagesStructs.Message{
				Action: "position",
			},
			Details: messagesStructs.PositionMessageDetails{
				X:     user.Location.X,
				Y:     user.Location.Y,
				Alpha: user.RotationAngle,
			},
		})
	}
	ticker.Stop()
}
