package apicore

import "sync"

import "doshelpv2/log"
import "doshelpv2/appctx"

import "golang.org/x/net/context"

type apiRequest struct {
	*ApiCore

	id uint64
	ctx context.Context
}
func (self *apiRequest) configure(ctx context.Context) ( *apiRequest, error ) {
	if self == nil { return nil,err_glob_InvalidSelf }
	if ctx == nil { return nil,err_glob_InvalidContext }

	self.ApiCore = ctx.Value(appctx.CTX_MOD_APICORE).(*ApiCore)
	self.slogger.W(log.LLEV_DBG, "Request submodule has been initialized and configured!")
	return self,nil
}


type request struct {
	ctx context.Context
}
type requestStore struct {
	sync.RWMutex
	store map[uint64]request
}
func (self *requestStore) configure() error {
	self.Lock(); defer self.Unlock()

	if self.store != nil { return err_ReqStore_StoreAlreadyDefined }
	self.store = make(map[uint64]request)
	return nil
}
func (self *requestStore) free() {	// Call it in DeInitModule 
	self.Lock()
	for i := range self.store { delete(self.store, i) }
	self.Unlock()
}
func (self *requestStore) put( id uint64, req request ) {
	self.Lock(); self.store[id] = req; self.Unlock()
}
func (self *requestStore) get( id uint64 ) ( *request, bool ) {
	self.RLock(); req, ok := self.store[id]; self.RUnlock()
	return &req, ok
}
func (self *requestStore) addRequest( ctx context.Context ) ( uint64, error ) {
	if self.store == nil { return uint64(0),err_ReqStore_StoreIsUndefined }
	if ctx == nil { return uint64(0),err_glob_InvalidInputData }

	var nextid uint64
	var req *request = new(request)

	nextid, e := self.getNextFreeID(); if e != nil { return uint64(0),e }
	req.ctx = ctx
	self.put(nextid, *req)
	if _, ok := self.get(nextid); ok == false {
		return uint64(0),err_ReqStore_ReqPutError
	} else { return nextid,nil }
	return uint64(0),nil
}
func (self *requestStore) getNextFreeID() ( uint64, error ) { return uint64(0),nil } // TODO
//func (self *requestStore) getLastUsedID() {}
//func (self *requestStore) dbInsert() {}
//func (self *requestStore) dbSync() {}
