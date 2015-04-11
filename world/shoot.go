package world

import (
	"CaribbeanWarServer/intmath"
	"CaribbeanWarServer/messagesStructs"
	"CaribbeanWarServer/rtree"
	"CaribbeanWarServer/structs"
	"time"
)

func (self *storage) shoot(user *structs.User, ch chan *messagesStructs.ShootIncome) {
	defer func() {
		sendErrorMessage(user, recover())
	}()
	for {
		if incomeMessage, ok := <-ch; ok {
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
			core := structs.NewCore(details.Location, user.RotationAngle, details.Angle, details.Direction, user.ID)
			user.Unlock()
			go self.updateCore(core, user)
			user.GetConn().WriteJSON(message)
		} else {
			return
		}
	}
}

func (self *storage) updateCore(core *structs.Core, user *structs.User) {
	timer := time.NewTicker(10 * time.Millisecond)
	defer func() {
		timer.Stop()
		sendErrorMessage(user, recover())
	}()
	for core.OverWater() {
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
			looser.SelectedShip.CurrentHP = uint16(intmath.Max(int(looser.SelectedShip.CurrentHP)-87, 0))
			if looser.SelectedShip.CurrentHP == 0 {
				message := messagesStructs.Dead{
					Action: "death",
					Details: messagesStructs.DeadDetails{
						ID:       looser.ID,
						Location: looser.Location,
						Rotation: looser.RotationAngle,
					},
				}
				looser.GetConn().WriteJSON(message)
				for _, neigbour := range looser.NearestUsers {
					neigbour.Conn.WriteJSON(message)
				}
			}
			looser.Unlock()
			return
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
