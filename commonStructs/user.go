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

type NearestUsers []NearestUser

func (self *NearestUsers) Len() int {
	return len(*self)
}

func (self *NearestUsers) Less(i, j int) bool {
	return (*self)[i].ID < (*self)[j].ID
}

func (self *NearestUsers) Swap(i, j int) {
	(*self)[i], (*self)[j] = (*self)[j], (*self)[i]
}
