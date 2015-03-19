package harbor

import (
	"CaribbeanWarServer/structs"
	"CaribbeanWarServer/world"
	"errors"
	"strconv"
	"sync"
)

type harborStruct struct {
	harbor []structs.User
	sync.Mutex
}

var harbor *harborStruct

func init() {
	harbor = &harborStruct{}
	harbor.harbor = make([]structs.User, 0, 30)
	world.InitHarbor(harbor.Add)
}

func GetHarbor() *harborStruct {
	return harbor
}

func (self *harborStruct) Add(user *structs.User) error {
	self.Lock()
	defer self.Unlock()
	if exist := self.exist(user.ID); exist {
		return errors.New("User exist")
	}
	self.harbor = append(self.harbor, *user)
	go self.waitForShipSelection(user)
	return nil
}

func (self *harborStruct) Remove(id uint) {
	self.Lock()
	defer self.Unlock()
	if index, err := self.indexOf(id); err == nil {
		self.harbor = append(self.harbor[:index], self.harbor[index+1:]...)
	}
}

func (self *harborStruct) exist(id uint) bool {
	_, err := self.indexOf(id)
	return err == nil
}

func (self *harborStruct) indexOf(id uint) (int, error) {
	for key, value := range self.harbor {
		if value.ID == id {
			return key, nil
		}
	}
	return 0, errors.New("Not found")
}

func (self *harborStruct) waitForShipSelection(user *structs.User) {
	defer func() {
		if err := recover(); err != nil {
			var message = "Something Wrong: Harbor: waitForShip"
			switch tmp := err.(type) {
			case error:
				message = tmp.Error()
			case string:
				message = tmp
			}
			self.sendErrorMessage(user.GetConn(), message)
			user.GetConn().Close()
			self.Remove(user.ID)
		}
	}()
	var data structs.Message
	if err := user.GetConn().ReadJSON(&data); err != nil {
		panic(err)
	}
	if data.Action != "enterWorld" {
		panic("Unrecognized action")
	}
	shipId := parseShipId(data.Details)
	selectedShip := findShipById(user.Ships, shipId)
	if selectedShip == nil {
		panic("Unrecognized ship")
	}
	user.SelectedShip = selectedShip
	world.Add(user)
	user.Lock()
	user.GetConn().WriteJSON(structs.Message{"enterWorld", map[string]interface{}{
		"success":      true,
		"nearestUsers": user.NearestUsers,
		"shipInfo":     user.SelectedShip,
		"location":     user.Location,
	}})
	user.Unlock()
	self.Remove(user.ID)
}

func (self *harborStruct) sendErrorMessage(conn *structs.Connection, err string) {
	conn.WriteJSON(structs.ErrorMessage(err))
}

func parseShipId(data map[string]interface{}) uint {
	switch tmp := data["shipId"].(type) {
	case string:
		i64, _ := strconv.ParseUint(tmp, 10, 0)
		return uint(i64)
	case float64:
		return uint(tmp)
	}
	return 0
}

func findShipById(ships []structs.Ship, shipId uint) *structs.Ship {
	for _, value := range ships {
		if value.ID == shipId {
			return &value
		}
	}
	return nil
}
