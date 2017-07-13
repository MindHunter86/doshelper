package p2p

import "golucky/system/util"

//import "github.com/nictuku/dht"
import "golang.org/x/net/context"
import "github.com/sirupsen/logrus"

type P2PService struct {
	log *logrus.Logger
	service_id uint8
	service_name string

	dht_infohash string
	dht_listenport uint16
}
func (self *P2PService) Configure(ctx context.Context) (uint8, util.Service, error) {
	if self == nil { return util.SERVICE_PTR_P2P,nil,util.Err_Glob_InvalidSelf }
	if ctx == nil { return util.SERVICE_PTR_P2P,nil,util.Err_Glob_InvalidContext }

	self.log = ctx.Value(util.CTX_MAIN_LOGGER).(*logrus.Logger)
	self.service_id = util.SERVICE_PTR_P2P
	self.service_name = util.SERVICE_PTR[self.service_id]
	self.dht_infohash = ctx.Value(util.CTX_MAIN_CONFIG).(*util.AppConfig).P2PInfoHash
	self.dht_listenport = ctx.Value(util.CTX_MAIN_CONFIG).(*util.AppConfig).P2PDhtLstnPort

	self.log.Debugln("Service " + self.service_name + " input config:")
	self.log.Debugln(self.dht_infohash)
	self.log.Debugln(self.dht_listenport)

	self.log.Debugln("Service " + self.service_name + " has been successfully configured! Service ready to run.")
	return self.service_id,self,nil
}
func (self *P2PService) Destroy() error {
	return nil
}
func (self *P2PService) Start() error {
	return nil
}
func (self *P2PService) Stop() error {
	return nil
}
//
//
//
//package main
//
//import "net/http"
//
//
//import "golucky/system/util"
//
//import "github.com/nictuku/dht"
//import "golang.org/x/net/context"
//import "github.com/sirupsen/logrus"
//
//
//func main() {
//
//	// example from https://github.com/meshbird/meshbird/blob/41a032c55c07b4b609cfa84d8ea365d4c58a8d59/common/discovery_dht.go:
//	var count int = 0
//
//	// dhtHash, e := dht.DecodeInfoHash("deca7a89a1dbdc4b213de1c0d5351e92582f31fb"); if e != nil {
//	dhtHash, e := dht.DecodeInfoHash("3D475467EA8A9723F0A8349503B77214945894AD"); if e != nil {
//		log.Fatalln(e)
//	}
//
//	dhtConfig := dht.NewConfig()
//	dhtConfig.Port = 9001
//	dhtConfig.DHTRouters = "router.utorrent.com:6881,router.magnets.im:6881,router.bittorrent.com:6881,dht.transmissionbt.com:6881,dht.aelitis.com:6881,router.bitcomet.com:6881"
//	dhtConfig.MaxNodes = 100000
//	dhtConfig.NumTargetPeers = 100
//	dhtConfig.RateLimit = 1000
//	dhtConfig.MaxInfoHashes = 10240
//	dhtConfig.MaxInfoHashPeers = 1024
//	dhtConfig.ClientPerMinuteLimit = 300
//
//
//	dhtInstance, e := dht.New(dhtConfig); if e != nil {log.Fatalln(e)}
//	if e = dhtInstance.Start(); e != nil {log.Fatalln(e)}
//
//	go func(instance *dht.DHT, infohash string) {
//		log.Println("SEND BROADCAST REQUEST")
//		instance.PeersRequest(infohash, true)
//	}(dhtInstance, string(dhtHash))
//
//	go http.ListenAndServe(":9011", nil)
//
//	for {
//		select {
//		case i := <-dhtInstance.PeersRequestResults:
//			count++
//			for hash,peers := range i {
//				// hash - dht.InfoHash
//				// peers - string slice
//				// var slen = len(peers)
//				// log.Printf("ID: %d; Hash: %s; Slice length: %d", count, hash, slen)
//				for _,encip := range peers {
//					log.Printf("ID: %d; Hash: %s; IP address: %s", count, hash, dht.DecodePeerAddress(encip))
//				}
////				for i,z := range v {
////					log.Println("ID:", count, "; ADDR:", dht.DecodePeerAddress(z))
////
////				}
//			}
////			for k,peers := range i {
////				for k2,peer := range peers {
////					log.Println(dht.DecodePeerAddress(peer))
////					if peerHash, e := dht.DecodeInfoHash(peer); e == nil {
////						log.Println(peerHash.String())
////					} else { log.Println("ERR:", e) }
////					log.Println("DEBUG: k-", k, "; k2-", k2)
////				}
////			}
//		}
//	}
//}
//
//
///*
//
//	router.utorrent.com:6881
//	router.magnets.im:6881
//	router.bittorrent.com:6881
//	dht.transmissionbt.com:6881
//	router.bitcomet.com:6881
//	dht.aelitis.com:6881
//
//2017/07/08 21:52:41 167.114.232.119:9000
//2017/07/08 21:52:41 ERR: DecodeInfoHash: expected InfoHash len=20, got 0
//2017/07/08 21:52:41 DEBUG: k- deca7a89a1dbdc4b213de1c0d5351e92582f31fb ; k2- 0
//
//*/
