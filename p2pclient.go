// Copyright 2019 New Vector Ltd
//
// This file is part of libp2p-proxy.
//
// libp2p-proxy is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// coap-proxy is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the Apache License v2
// along with libp2p-proxy.  If not, see <https://www.gnu.org/licenses/>.

package main

import (

	"context"
	"crypto/rand"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"time"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"

	ds "github.com/ipfs/go-datastore"
	dsync "github.com/ipfs/go-datastore/sync"

	dht "github.com/libp2p/go-libp2p-kad-dht"
	rhost "github.com/libp2p/go-libp2p/p2p/host/routed"
	ma "github.com/multiformats/go-multiaddr"
)



func MakeRoutedHost(
	listenPort int,
	privKey crypto.PrivKey,
	bootstrapPeers []peer.AddrInfo,
) (host.Host, error) {

	// TODO: review whether this is the correct thing to do with contexts.
	ctx := context.Background()

	basicHost, err := libp2p.New(
		ctx,
		libp2p.Identity(privKey),
		libp2p.ListenAddrStrings(fmt.Sprintf("/ip4/0.0.0.0/tcp/%s", strconv.Itoa(p2pPort))),
		libp2p.DefaultTransports,
		libp2p.DefaultMuxers,
		libp2p.DefaultSecurity,
		libp2p.NATPortMap(),
	)

	if err != nil {
		return nil, err
	}

	log.Println("Started host with id %s", basicHost.ID())

	// Construct a datastore (needed by the DHT). This is just a simple,
	// in-memory thread-safe datastore.
	dstore := dsync.MutexWrap(ds.NewMapDatastore())

	// Make the DHT
	dht_router := dht.NewDHT(ctx, basicHost, dstore)

	// Make the routed host
	routedHost := rhost.Wrap(basicHost, dht_router)

	// connect to the chosen ipfs nodes
	err = bootstrapConnect(ctx, routedHost, bootstrapPeers)

	if err != nil {
		return nil, err
	}

	bootstrapConf := dht.BootstrapConfig{
	    Queries: 2,
	    Period: time.Duration(30 * time.Second),
	    Timeout: time.Duration(30 * time.Second),
	}

	// Bootstrap the host
	err = dht_router.BootstrapWithConfig(ctx, bootstrapConf)

	if err != nil {
		return nil, err
	}

	// Build host multiaddress
	hostAddr, _ := ma.NewMultiaddr(fmt.Sprintf("/ipfs/%s", routedHost.ID().Pretty()))

	// Now we can build a full multiaddress to reach this host
	// by encapsulating both addresses:
	addrs := routedHost.Addrs()
	for _, addr := range addrs {
		log.Println(addr.Encapsulate(hostAddr))
	}

	log.Println("I can be reached at:")
	for _, addr := range addrs {
		log.Println(addr.Encapsulate(hostAddr))
	}

	return routedHost, nil
}

// Retrieve the identity (the private key) for the p2p peer
// Tries to retrieve the identity from the specified file.
// If the file doesn't exist it will try to create it.
// If the content of the file is not a valid encoded p2p
// PrivKey it will exit.
func GetLibp2pPrivKey(identityFile string) crypto.PrivKey {
	encodedPrivateKey, err := ioutil.ReadFile(identityFile)

	var privKey crypto.PrivKey

	if err != nil {
		log.Printf("Couldn't find an identity in %s, creating new idenitity.\n", identityFile)
		privKey = CreatePeerIdentity(identityFile)
		log.Println("Identity created")
	} else {
		handler := func(err error) {
			if err != nil {
				log.Fatalf("The identity in %s is invalid. Shutting down proxy.", identityFile)
			}
		}

		decodedPrivateKey, err := crypto.ConfigDecodeKey(string(encodedPrivateKey))

		handler(err)

		privKey, err = crypto.UnmarshalPrivateKey(decodedPrivateKey)

		handler(err)
	}

	return privKey
}

// Creates a p2p identity in the form of a private key.
func CreatePeerIdentity(identityFile string) crypto.PrivKey {
	r := rand.Reader

	privKey, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, r)

	if err != nil {
		panic(err)
	}

	privKeyBytes, err := crypto.MarshalPrivateKey(privKey)

	if err != nil {
		panic(err)
	}

	encodedPrivateKey := crypto.ConfigEncodeKey(privKeyBytes)

	// Only the current user has access to the identity file
	err = ioutil.WriteFile(identityFile, []byte(encodedPrivateKey), 0600)

	if err != nil {
		log.Fatalf("Couldn't write the idenity file to %s", identityFile)
	}

	return privKey
}
