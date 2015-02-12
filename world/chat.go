package world

import (
	"CaribbeanWarServer/rtree"
	"CaribbeanWarServer/structs"
)

func (self *storage) chat(message *map[string]interface{}) {
	self.Lock()
	self.ocean.Each(func(s *rtree.Spatial) {
		(*s).(*structs.User).GetConn().WriteJSON(message)
	})
	self.Unlock()
}
