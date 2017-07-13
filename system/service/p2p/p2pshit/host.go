package main

import (
	"context"
	"crypto/rand"
	"fmt"
//	"encoding/hex"

	crypto "github.com/libp2p/go-libp2p-crypto"
	peer "github.com/libp2p/go-libp2p-peer"
	pstore "github.com/libp2p/go-libp2p-peerstore"
	swarm "github.com/libp2p/go-libp2p-swarm"
	bhost "github.com/libp2p/go-libp2p/p2p/host/basic"
	ma "github.com/multiformats/go-multiaddr"
	ds "github.com/ipfs/go-datastore"
	"github.com/libp2p/go-libp2p-kad-dht"
//	"github.com/multiformats/go-multihash"
	"github.com/ipfs/go-cid"
)

func main() {
	// Generate an identity keypair using go's cryptographic randomness source
	priv, pub, err := crypto.GenerateEd25519Key(rand.Reader)
	if err != nil {
		panic(err)
	}

	// A peers ID is the hash of its public key
	pid, err := peer.IDFromPublicKey(pub)
	if err != nil {
		panic(err)
	}

	// We've created the identity, now we need to store it.
	// A peerstore holds information about peers, including your own
	ps := pstore.NewPeerstore()
	ps.AddPrivKey(pid, priv)
	ps.AddPubKey(pid, pub)

	maddr, err := ma.NewMultiaddr("/ip4/0.0.0.0/tcp/9000")
	if err != nil {
		panic(err)
	}

	// Make a context to govern the lifespan of the swarm
	ctx := context.Background()

	// Put all this together
	netw, err := swarm.NewNetwork(ctx, []ma.Multiaddr{maddr}, pid, ps, nil)
	if err != nil {
		panic(err)
	}

	myhost := bhost.New(netw)
	fmt.Printf("Hello World, my hosts ID is %s\n", myhost.ID())

// KAD TESTING:
	// testds := ds.Datastore{}
//	fmt.Println(testds)
	testds := ds.NewNullDatastore()

	dhtcli := dht.NewDHT(ctx, myhost, testds)
	fmt.Println( dhtcli.Close() )

	e := dhtcli.Bootstrap(dhtcli.Context()); if e != nil {fmt.Println(e)}

	// cid.Cast - <version><codec-type><multihash>
	cidConfig := []byte{1,1}

//	hexString, e := hex.DecodeString("deca7a89a1dbdc4b213de1c0d5351e92582f31fb")
//	if e != nil {fmt.Println(e)}
//	encKeyString, e := multihash.EncodeName(hexString, "sha1")
//	if e != nil {fmt.Println(e)}

	//multiDec, e := multihash.Decode(encKeyString); if e != nil {fmt.Println(e)}
	///multiHex := hex.EncodeToString(encKeyString)
	cidConfig = append([]byte{1,1,0x11, byte(len("0beec7b5ea3f0fdbc95d0dd47f3c5bc275da8a33"))}, "0beec7b5ea3f0fdbc95d0dd47f3c5bc275da8a33"...)
	//cidConfig = append(cidConfig, "deca7a89a1dbdc4b213de1c0d5351e92582f31fb"...)
	//cidConfig = append(cidConfig, "QmbFMke1KXqnYyBBWxB74N4c5SBnJMVAiMNRcGu6x1AwQH"...)
	//cidConfig = append(cidConfig, t1)// "0beec7b5ea3f0fdbc95d0dd47f3c5bc275da8a33"...)
	gocid,e := cid.Cast(cidConfig); if e != nil {fmt.Println(e)}

	fmt.Println(gocid.String())
	e = dhtcli.Provide(dhtcli.Context(), gocid, true); if e != nil {fmt.Println(e)}

	fmt.Println("kad")
//	pstore, e := dhtcli.FindPeersConnectedToPeer(ctx, pid); if e !=nil {fmt.Errorf("FUCK,", e)}

//	for {
//		select {
//		case pinfo := <-pstore:
//			fmt.Println("Found:", pinfo.ID.String())
//		}
//	}
	kadtest, e := dhtcli.GetClosestPeers(dhtcli.Context(), "deca7a89a1dbdc4b213de1c0d5351e92582f31fb"); if e != nil {fmt.Println("fuck,", e)}

	for {
		select {
		case i:=<-kadtest:
			fmt.Println("FOUND!", i.String())
		}
	}
}
