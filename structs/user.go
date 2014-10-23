package structs

import (
	"CaribbeanWarServer/rtree"
	"github.com/gorilla/websocket"
)

type User struct {
	ID           uint
	Email        string
	Nick         string
	Cash         uint
	Conn         *websocket.Conn
	Location     Point
	Ships        []Ship
	SelectedShip *Ship
	NearestUsers []*User
	Rotation     float64
	Speed        float64
}

func (self *User) Bounds() *rtree.Rect {
	bound, _ := rtree.NewRect(rtree.Point{self.Location.X, self.Location.Y}, []float64{1, 1})
	return bound
}
