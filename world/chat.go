package world

import (
	"CaribbeanWarServer/rtree"
	"CaribbeanWarServer/structs"
)

func (self *storage) chat(message *structs.Message) {
	self.ocean.Each(func(s *rtree.Spatial) {
		(*s).(*structs.User).GetConn().WriteJSON(message)
	})
}
