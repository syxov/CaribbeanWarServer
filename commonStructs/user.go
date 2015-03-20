package commonStructs

import (
	"CaribbeanWarServer/point"
)

type NearestUser struct {
	ID            uint          `json:"id"`
	Conn          *Connection   `json:"-"`
	Ship          *Ship         `json:"ship"`
	Nick          string        `json:"nick"`
	Location      point.Point2D `json:"location"`
	RotationAngle float64       `json:"alpha"`
}
