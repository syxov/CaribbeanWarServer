package world

import (
	"CaribbeanWarServer/messagesStructs"
	"CaribbeanWarServer/rtree"
	"CaribbeanWarServer/structs"
)

func (self *storage) chat(message *messagesStructs.Message) {
	go self.ocean.Each(func(s *rtree.Spatial) {
		(*s).(*structs.User).GetConn().WriteJSON(message)
	})
}
