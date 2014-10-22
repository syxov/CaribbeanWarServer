package world

import (
	"CaribbeanWarServer/rtree"
	"CaribbeanWarServer/structs"
)

func (self *storage) chat(data map[string]interface{}) {
	message := map[string]interface{}{
		"action": "chat",
		"details": map[string]interface{}{
			"sender":  data["sender"],
			"message": data["message"],
		},
	}
	self.Lock()
	defer self.Unlock()
	self.ocean.Each(func(s *rtree.Spatial) {
		(*s).(*structs.User).Conn.WriteJSON(message)
	})
}
