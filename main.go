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
	"flag"
	"log"
	"net/http"

	"github.com/gorilla/mux"

//	gologging "github.com/whyrusleeping/go-logging"
//	golog "github.com/ipfs/go-log"
)

var (
	httpTarget         = *flag.String("http-target", "http://127.0.0.1:8008", "The HTTP host+port the requests to this peer are sent to")
	httpPort           = *flag.Int("http-port", 8999, "The HTTP port to listen on")
	p2pPort            = *flag.Int("p2p-port", 8998, "The port libp2p listens on for p2p communication.")
	identityFile       = *flag.String("identity-file", "./identity", "A file containig the libp2p peer's private key.")
	bootstrapPeersFile = *flag.String("peers-file", "./bootstrap", "A file containing ipfs peers used for discovering other peers via id alone.")
)

func init() {
	flag.Parse()
}

func main() {
	// LibP2P code uses golog to log messages. They log with different
	// string IDs (i.e. "swarm"). We can control the verbosity level for
	// all loggers with:
//	golog.SetAllLoggers(gologging.INFO) // Change to DEBUG for extra info

	privKey := GetLibp2pPrivKey(identityFile)

	host, err := MakeRoutedHost(p2pPort, privKey, IPFS_PEERS)

	proxyService := NewProxyService(host)

	if err != nil {
		panic(err)
		log.Fatal("Couldn't start libp2p host.")
	}

	// Set up the http handlers
	r := mux.NewRouter()
	r.HandleFunc("/{.*}", proxyService.ServeHTTP).Host("{identity:.*}.matrixp2p")
	r.HandleFunc("/{.*}", proxyService.ServeHTTPS)

	http.ListenAndServe(":7667", r)


	select {}
}

