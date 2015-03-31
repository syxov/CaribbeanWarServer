package world

import (
	"CaribbeanWarServer/messagesStructs"
	"CaribbeanWarServer/rtree"
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
	core := structs.NewCore(details.Location, details.Angle, user.RotationAngle, details.Direction, user.ID)
	user.Unlock()
	go self.updateCore(core, user)
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
		spatials := self.ocean.SearchIntersectWithLimit(1, core.GetBounds(), func(spat *rtree.Spatial) bool {
			return (*spat).(*structs.User).ID != user.ID
		})
		if len(spatials) != 0 {
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
			break
		}
	}
	user.GetConn().WriteJSON(map[string]interface{}{
		"action": "miss",
		"details": map[string]interface{}{
			"position": core.CurrentPosition,
		},
	})
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
