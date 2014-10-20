package world

import (
	"CaribbeanWarServer/structs"
	"github.com/dhconnelly/rtreego"
	"github.com/gorilla/websocket"
)

type node struct {
	ID       uint
	Conn     *websocket.Conn
	Location *structs.Point
	bound    *rtreego.Rect
}

func (self *node) Bounds() *rtreego.Rect {
	if self.bound == nil {
		self.bound, _ = rtreego.NewRect(rtreego.Point{self.Location.X, self.Location.Y}, []float64{1, 1})
	}
	return self.bound
}
