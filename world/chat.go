package world

import (
	"CaribbeanWarServer/quadtree"
	"CaribbeanWarServer/structs"
)

func (self *storage) chat(message *map[string]interface{}) {
	self.Lock()
	defer self.Unlock()
	self.ocean.Each(func(s *quadtree.Point) {
		s.Data().(*structs.User).GetConn().WriteJSON(message)
	})
}
