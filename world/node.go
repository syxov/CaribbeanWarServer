package world

import (
	"CaribbeanWarServer/structs"
	"github.com/dhconnelly/rtreego"
)

type node struct {
	ID       uint
	Location *structs.Point
	bound    *rtreego.Rect
}

func (self *node) Bounds() *rtreego.Rect {
	if self.bound == nil {
		self.bound, _ = rtreego.NewRect(rtreego.Point{self.Location.X, self.Location.Y}, []float64{1, 1})
	}
	return self.bound
}
