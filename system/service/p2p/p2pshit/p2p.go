package main

import "log"
import "context"
import "crypto/rand"

import crypto "github.com/libp2p/go-libp2p-crypto"
import peer "github.com/libp2p/go-libp2p-peer"
import pstore "github.com/libp2p/go-libp2p-peerstore"
import ma "github.com/multiformats/go-multiaddr"
import swarm "github.com/libp2p/go-libp2p-swarm"
import bhost "github.com/libp2p/go-libp2p/p2p/host/basic"


func main() {
	log.Println("Fuck my feels, baby")
	log.Println("tes")


	prv, pub, e := crypto.GenerateEd25519Key(rand.Reader)
	if e != nil { log.Fatalln(e) }

	pid, e := peer.IDFromPublicKey(pub)
	if e != nil { log.Fatalln(e) }

	ps := pstore.NewPeerstore()
	ps.AddPrivKey(pid, prv)
	ps.AddPubKey(pid, pub)

	maddr, e := ma.NewMultiaddr("/ip4/0.0.0.0/tcp/9000")
	if e != nil { log.Fatalln(e) }

	ctx := context.Background()
	netw, e := swarm.NewNetwork(ctx, []ma.Multiaddr{maddr}, pid, ps, nil)
	if e != nil { log.Fatalln(e) }

	myhost := bhost.New(netw)
	log.Println("OOK", myhost.ID().String())
}
