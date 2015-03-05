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
			var message string
			switch tmp := err.(type) {
			case error:
				message = tmp.Error()
			default:
				message = "Something Wrong: Harbor: waitForShip"
			}
			self.sendErrorMessage(user.GetConn(), errors.New(message))
			user.GetConn().Close()
			self.Remove(user.ID)
		}
	}()
	var dataI map[string]interface{}
	if err := user.GetConn().ReadJSON(&dataI); err == nil {
		action := dataI["action"].(string)
		if action == "enterWorld" {
			var id uint
			switch tmp := dataI["details"].(map[string]interface{})["shipId"].(type) {
			case string:
				i64, _ := strconv.ParseUint(tmp, 10, 0)
				id = uint(i64)
			case float64:
				id = uint(tmp)
			}
			for _, value := range user.Ships {
				if value.ID == id {
					user.SelectedShip = &value
					world.Add(user)
					user.Lock()
					user.GetConn().WriteJSON(map[string]interface{}{
						"action": "enterWorld",
						"details": map[string]interface{}{
							"success":      true,
							"nearestUsers": user.NearestUsers,
							"shipInfo":     user.SelectedShip,
							"location":     user.Location,
						},
					})
					user.Unlock()
					self.Remove(user.ID)
					return
				}
			}
			panic(errors.New("Unrecognized ship"))
		} else {
			panic(errors.New("Unrecognized action"))
		}
	} else {
		panic(err)
	}
}

func (self *harborStruct) sendErrorMessage(conn *structs.Connection, err error) {
	conn.WriteJSON(map[string]interface{}{
		"action": "fuckup",
		"details": map[string]string{
			"message": err.Error(),
		},
	})
}
