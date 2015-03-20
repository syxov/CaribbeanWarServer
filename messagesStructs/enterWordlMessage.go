package messagesStructs

import (
	"CaribbeanWarServer/commonStructs"
	"CaribbeanWarServer/point"
)

type EnterWorld struct {
	Action  string            `json:"action"`
	Details EnterWorldDetails `json:"details"`
}

type EnterWorldDetails struct {
	Success      bool                        `json:"success"`
	NearestUsers []commonStructs.NearestUser `json:"nearestUsers"`
	ShipInfo     commonStructs.Ship          `json:"shipInfo"`
	Location     point.Point2D               `json:"location"`
}
