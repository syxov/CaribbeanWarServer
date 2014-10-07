package world

func (self *WorldStruct) shop(data map[string]interface{}) {
	requestType := data["details"].(map[string]string)["request"]
	id := data["details"].(map[string]uint)["param"]
	if requestType == "load" {
		self.DbConn.GetShopItems(id)
	} else {

	}
}
