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
    "flag"
    "io/ioutil"
    "log"

    "github.com/libp2p/go-libp2p"
    "github.com/libp2p/go-libp2p-core/crypto"
)

var (
    httpTarget   = flag.String("http-target", "http://127.0.0.1:8008", "The HTTP host+port the requests to this peer are sent to")
    httpPort     = flag.String("http-port", "8888", "The HTTP port to listen on")
    identityFile = flag.String("identity-file", "./identity", "A file containig the libp2p peer's private key.")
)

func init() {
    flag.Parse()
}

func main() {

    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    privKey := getLibp2pPrivKey(*identityFile)

    h, err := libp2p.New(
	ctx,
	libp2p.Identity(privKey),
    )

    if err != nil {
	log.Fatal("Couldn't start libp2p host.")
    }

    log.Printf("Started host with id %s", h.ID())

}

// Retrieve the identity (the private key) for the p2p peer
// Tries to retrieve the identity from the specified file.
// If the file doesn't exist it will try to create it.
// If the content of the file is not a valid encoded p2p
// PrivKey it will exit.
func getLibp2pPrivKey(identityFile string) crypto.PrivKey {
    encodedPrivateKey, err := ioutil.ReadFile(identityFile)

    var privKey crypto.PrivKey

    if err != nil {
	log.Printf("Couldn't find an identity in %s, creating new idenitity.\n", identityFile)
	privKey = createPeerIdentity(identityFile)
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

func createPeerIdentity(identityFile string) crypto.PrivKey {
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


