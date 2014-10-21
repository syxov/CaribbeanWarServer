package harbor

import (
	"CaribbeanWarServer/structs"
	"CaribbeanWarServer/world"
	"errors"
	"github.com/gorilla/websocket"
	"sync"
	"time"
)

type HarborStruct struct {
	harbor []structs.User
	sync.Mutex
}

func (self *HarborStruct) Add(user *structs.User) error {
	self.Lock()
	defer self.Unlock()
	if cap(self.harbor) == 0 {
		self.harbor = make([]structs.User, 0, 30)
	}
	if exist := self.exist(user.ID); exist {
		return errors.New("User exist")
	}
	self.harbor = append(self.harbor, *user)
	go self.waitForShipSelection(user)
	go ping(user.Conn)
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
				message = "Somethign Wrong: Harbor: waitForShip"
			}
			self.sendErrorMessage(user.Conn, errors.New(message))
			user.Conn.Close()
			self.Remove(user.ID)
		}
	}()
	var dataI interface{}
	if err := user.Conn.ReadJSON(&dataI); err == nil {
		data := dataI.(map[string]interface{})
		action := data["Action"].(string)
		if action == "shipSelect" {
			id := uint(data["Details"].(map[string]interface{})["ShipId"].(float64))
			for _, value := range user.Ships {
				if value.ID == id {
					user.SelectedShip = &value
					user.Conn.WriteJSON(map[string]interface{}{
						"Action": "shipSelect",
						"Details": map[string]bool{
							"Success": true,
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
		"Action": "fuckup",
		"Details": map[string]string{
			"Message": err.Error(),
		},
	})
}

func ping(conn *websocket.Conn) {
	for {
		time.Sleep(13 * time.Second)
		if err := conn.WriteMessage(websocket.TextMessage, []byte("{}")); err != nil {
			return
		}
	}
}
