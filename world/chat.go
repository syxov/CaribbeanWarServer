package world

import (
	"CaribbeanWarServer/messagesStructs"
	"CaribbeanWarServer/rtree"
	"CaribbeanWarServer/structs"
)

func (self *storage) chat(ch chan *messagesStructs.Message) {
	for {
		if message, ok := <-ch; ok {
			self.ocean.Each(func(s *rtree.Spatial) {
				(*s).(*structs.User).GetConn().WriteJSON(message)
			})
		} else {
			return
		}
	}
}
