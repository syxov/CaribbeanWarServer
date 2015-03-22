package world

import (
	"CaribbeanWarServer/messagesStructs"
	"CaribbeanWarServer/structs"
	"time"
)

func (self *storage) shoot(user *structs.User, incomeMessage messagesStructs.ShootIncome) {
	defer func() {
		sendErrorMessage(user, recover())
	}()
	details := incomeMessage.Details
	user.Lock()
	message := messagesStructs.ShootOutcome{
		Action: "shoot",
		Details: messagesStructs.ShootOutcomeDetails{
			ID:        user.ID,
			Alpha:     user.RotationAngle,
			Angle:     details.Angle,
			Location:  details.Location,
			Direction: details.Direction,
		},
	}
	for _, neigbour := range user.NearestUsers {
		neigbour.Conn.WriteJSON(message)
	}
	core := structs.NewCore(&details.Location, user.RotationAngle, details.Angle, details.Direction, user.ID)
	go self.updateCore(core, user)
	user.Unlock()
	user.GetConn().WriteJSON(message)
}

func (self *storage) updateCore(core *structs.Core, user *structs.User) {
	defer func() {
		sendErrorMessage(user, recover())
	}()
	timer := time.NewTicker(10 * time.Millisecond)
	defer timer.Stop()
	for !core.UnderWater() {
		now := time.Now().UnixNano()
		<-timer.C
		core.UpdatePosition(float64(time.Now().UnixNano()-now) / float64(time.Second))
		spatials := self.ocean.SearchIntersectWithLimit(1, core.GetBounds())
		if len(spatials) == 1 {
			looser := spatials[0].(*structs.User)
			message := messagesStructs.Hit{
				Action: "hit",
				Details: messagesStructs.HitDetails{
					ID:       looser.ID,
					Location: core.CurrentPosition,
					Damage:   87,
				},
			}
			looser.Lock()
			looser.GetConn().WriteJSON(message)
			for _, neigbour := range looser.NearestUsers {
				neigbour.Conn.WriteJSON(message)
			}
			looser.Unlock()
			return
		}
	}
}

func sendErrorMessage(user *structs.User, err interface{}) {
	if err != nil {
		var message = "Smth wrong: shoot"
		switch smth := err.(type) {
		case error:
			message = smth.Error()
		case string:
			message = smth
		}
		user.GetConn().WriteJSON(messagesStructs.ErrorMessage(message))
	}
}
