
package main

import (
	"net/http"
	"bufio"
	"log"

	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	ma "github.com/multiformats/go-multiaddr"
)

const Protocol = "/proxy/1.0.0"

type ProxyService struct {
    host host.Host
    dest peer.ID
    proxyAddr ma.Multiaddr
}

func NewProxyService(h host.Host) *ProxyService {
    h.SetStreamHandler(Protocol, ServeLibp2pProxy)

    log.Println("Proxy server is ready")
    log.Println("libp2p-peer addresses:")
    for _, a := range h.Addrs() {
	    log.Printf("%s/ipfs/%s\n", a, peer.IDB58Encode(h.ID()))
    }

    return &ProxyService {
	host: h,
    }
}

// Handles a stream over libp2p as though it were a proxy request.
func ServeLibp2pProxy(stream network.Stream) {
    defer stream.Close()

    buf := bufio.NewReader(stream)

    req, err := http.ReadRequest(buf)
    if err != nil {
	stream.Reset()
	log.Println(err)
    }
    defer req.Body.Close()

    req.URL.Scheme = "http"

    req.URL.Host = httpTarget

    outreq := new(http.Request)
    *outreq = *req

    log.Printf("Making request to %s\n.", req.URL)

    resp, err := http.DefaultTransport.RoundTrip(outreq)
    if err != nil {
	stream.Reset()
	log.Println(err)
	return
    }

    resp.Write(stream)
}
