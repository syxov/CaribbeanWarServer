package harbor

import (
	"CaribbeanWarServer/structs"
	"CaribbeanWarServer/world"
	"errors"
	"github.com/gorilla/websocket"
	"strconv"
	"sync"
)

type HarborStruct struct {
	harbor []structs.User
	sync.Mutex
}

var harbor HarborStruct

func init() {
	harbor.harbor = make([]structs.User, 0, 30)
	world.InitHarbor(harbor.add)
}

func Add(user *structs.User) error {
	return harbor.add(user)
}

func (self *HarborStruct) add(user *structs.User) error {
	self.Lock()
	defer self.Unlock()
	if exist := self.exist(user.ID); exist {
		return errors.New("User exist")
	}
	self.harbor = append(self.harbor, *user)
	go self.waitForShipSelection(user)
	return nil
}

func (self *HarborStruct) Remove(id uint) {
	if index, err := self.indexOf(id); err == nil {
		self.Lock()
		defer self.Unlock()
		self.harbor = append(self.harbor[:index], self.harbor[index+1:]...)
	}
}

func (self *HarborStruct) exist(id uint) bool {
	_, err := self.indexOf(id)
	return err == nil
}

func (self *HarborStruct) indexOf(id uint) (int, error) {
	for key, value := range self.harbor {
		if value.ID == id {
			return key, nil
		}
	}
	return 0, errors.New("Not found")
}

func (self *HarborStruct) waitForShipSelection(user *structs.User) {
	defer func() {
		if err := recover(); err != nil {
			var message string
			switch err.(type) {
			case error:
				message = err.(error).Error()
			default:
				message = "Something Wrong: Harbor: waitForShip"
			}
			self.sendErrorMessage(user.Conn, errors.New(message))
			user.Conn.Close()
			self.Remove(user.ID)
		}
	}()
	var dataI interface{}
	if err := user.Conn.ReadJSON(&dataI); err == nil {
		data := dataI.(map[string]interface{})
		action := data["action"].(string)
		if action == "enterWorld" {
			var id uint
			switch tmp := data["details"].(map[string]interface{})["shipId"].(type) {
			case string:
				i64, _ := strconv.ParseUint(tmp, 10, 0)
				id = uint(i64)
			case float64:
				id = uint(tmp)
			}
			for _, value := range user.Ships {
				if value.ID == id {
					user.SelectedShip = &value
					user.Conn.WriteJSON(map[string]interface{}{
						"action": "enterWorld",
						"details": map[string]bool{
							"success": true,
						},
					})
					world.Add(user)
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

func (self *HarborStruct) sendErrorMessage(conn *websocket.Conn, err error) {
	conn.WriteJSON(map[string]interface{}{
		"action": "fuckup",
		"details": map[string]string{
			"message": err.Error(),
		},
	})
}
