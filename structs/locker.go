package structs

import (
	"fmt"
	"github.com/zond/gotomic"
	"runtime/debug"
	"sync"
	"time"
)

type Locker struct {
	StackTrace *gotomic.List
	Time       *time.Timer
	ownMutex   sync.Mutex
	sync.Mutex
}

func (self *Locker) Lock() {
	self.ownMutex.Lock()
	if self.StackTrace == nil {
		self.StackTrace = gotomic.NewList()
	}
	self.StackTrace.Push(debug.Stack())
	self.Time = time.AfterFunc(10*time.Second, func() {
		self.ownMutex.Lock()
		self.StackTrace.Each(func(t gotomic.Thing) bool {
			fmt.Println(string(t.([]byte)))
			return false
		})
		fmt.Println("\n--------------------------------------\n")
		self.ownMutex.Unlock()
	})
	self.ownMutex.Unlock()
	self.Mutex.Lock()
}

func (self *Locker) Unlock() {
	self.ownMutex.Lock()
	self.Time.Stop()
	self.StackTrace.Pop()
	self.ownMutex.Unlock()
	self.Mutex.Unlock()
}
