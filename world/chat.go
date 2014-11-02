package world

import (
	"CaribbeanWarServer/rtree"
	"CaribbeanWarServer/structs"
)

func (self *storage) chat(message *interface{}) {
	self.Lock()
	self.ocean.Each(func(s *rtree.Spatial) {
		(*s).(*structs.User).Conn.WriteJSON(message)
	})
	self.Unlock()
}
