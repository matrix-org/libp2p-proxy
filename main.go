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
    "fmt"

//    "github.com/libp2p/go-libp2p"
)

func main() {
    httpTarget := flag.String("http-target", "http://127.0.0.1:8008", "The HTTP host+port the requests to this peer are sent to")
    httpPort := flag.String("http-port", "8888", "The HTTP port to listen on")
    identityFile := flag.String("identity-file", "./identity", "A file containig the libp2p peer's private key.")

    fmt.Printf("%s %s %s", *httpTarget, *httpPort, *identityFile)
}
