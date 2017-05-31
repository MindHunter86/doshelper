package apishield

import "sync"
import"time"

type guest struct {
	uuid string
	addr, uagent, origin, referer []byte
	first_access, last_access time.Time
}
func (self *guest) create( access_time time.Time) (*guest,error) {
	if len(access_time.String()) == 0 { return nil,err_mod_InvalidInputData }
	self.first_access = access_time
	return self,nil
}

// in submodule guest-store.go
type guestStore struct{
	sync.RWMutex
	store map[string]guest
}
func (self *guestStore) create() (*guestStore,error) {
	if self == nil { return nil,err_mod_InvalidSelf }
	if self.store != nil { return nil,err_gstore_StoreAlreadyDefined }

	self.Lock(); self.store = make(map[string]guest); self.Unlock()
	return self,nil
}
func (self *guestStore) destroy() error {
	if self == nil { return err_mod_InvalidSelf }
	if self.store == nil { return err_gstore_InvalidStore }

	self.Lock(); for key := range self.store { delete(self.store, key) }; self.Unlock()
	return nil
}
func (self *guestStore) get(key string) (*guest,error) {
	if self.store == nil { return nil,err_gstore_InvalidStore }
	if len(key) == 0 { return nil,err_mod_InvalidInputData }

	self.RLock(); t1, ok := self.store[key]; self.RUnlock()
	if ok == false { return nil,err_gstore_RecordNotFound }
	return &t1,nil
}
func (self *guestStore) put(key string, gst *guest) error {
	if self.store == nil { return err_gstore_InvalidStore }
	if gst == nil || len(key) == 0 { return err_mod_InvalidInputData }

	self.Lock(); self.store[key] = *gst ; self.Unlock()
	return nil
}
