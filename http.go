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
    "io"
    "net/http"
    "log"
    "bufio"

    "github.com/libp2p/go-libp2p-core/peer"
    "github.com/gorilla/mux"
	ma "github.com/multiformats/go-multiaddr"
)

// Takes a http request whose host is of the form *.matrixp2p
// and forwards it accross the libp2p connection.
func (p *ProxyService) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    addr:= "/ipfs/" + params["identity"]

    ipfsaddr, err := ma.NewMultiaddr(addr)
    if err != nil {
	log.Printf("Got a request for an invalid peer id %s", params["identity"])
	http.Error(w, err.Error(), http.StatusServiceUnavailable)
	return
    }

    pid, err := ipfsaddr.ValueForProtocol(ma.P_IPFS)
    if err != nil {
	log.Printf("Got a request for an invalid peer id %s", params["identity"])
	http.Error(w, err.Error(), http.StatusServiceUnavailable)
	return
    }

    peerID, err := peer.IDB58Decode(pid)
    if err != nil {
	log.Printf("Got a request for an invalid peer id %s", params["identity"])
	http.Error(w, err.Error(), http.StatusServiceUnavailable)
	return
    }


    log.Printf("Got a request for %s", peerID)

    stream, err := p.host.NewStream(context.Background(), peerID, Protocol)
    if err != nil {
	log.Println(err)
	http.Error(w, err.Error(), http.StatusInternalServerError)
	return
    }

    defer stream.Close()

    err = r.Write(stream)
    if err != nil {
	stream.Reset()
	log.Println(err)
	http.Error(w, err.Error(), http.StatusServiceUnavailable)
	return
    }

    // Read the response sent by the peer
    buf := bufio.NewReader(stream)
    resp, err := http.ReadResponse(buf, r)

    if err != nil {
	stream.Reset()
	log.Println(err)
	http.Error(w, err.Error(), http.StatusServiceUnavailable)
	return
    }

    for k, v := range resp.Header {
	for _, s := range v {
	    w.Header().Add(k, s)
	}
    }

    io.Copy(w, resp.Body)
    resp.Body.Close()
}

// Takes a https request and forwards it to the host via a
// tcp tunnel
func (p *ProxyService) ServeHTTPS(w http.ResponseWriter, r *http.Request) {

}
