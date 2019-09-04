The libp2p proxy can be put in front of any homeserver, whether that server is
intended to be run as a p2p peer or not. In this fashion we can bridge the gap
between HTTP transport homeservers and HTTP over libp2p transport homeservers.


The transport over libp2p has it's own authentication and encryption system.
This requires the homeserver to use HTTP without tls for `*.matrixp2p`
domains. This is a suggested option for homeservers according to the spec.
It's not required however, so you should check that your homeserver
implementation supports it. At the moment of writing this overview (Sep 2019)
only Synapse supports this.


If a request is to a domain without the `*.matrixp2p` tld libp2p-proxy will
act like a tcp proxy (much like a standard proxy's CONNECT method)


The structure of the network would look somewhat like this:


```
                                               +-----------------+                                     
                                               |                 |                                     
                                    +--------->|   homeserver    |<---------+                          
                                    |          |   example.org   |          |                          
                                    |          |                 |          |                          
                                    |          +-----------------+          |                          
                                    |                   |                   |                          
                                    |                   | HTTP              |                          
                                    |                   v                   |                          
                                    |            +--------------+           |                          
                                    |            |              |           |                          
                                    |            | libp2p-proxy |           |                          
                              HTTPS |            |              |           | HTTPS                         
                                    |            +--------------+           |                          
                                    |                 |   |                 |                          
                                    |                 |   |                 |                          
                                    |                 |   |                 |                          
                                    |          +------+   +------+          |                          
                                    |          |                 |          |                          
                                    |          |                 |          |                          
                                    |          |                 |          |                          
                                    |          |                 |          |                          
                                    |          v                 v          |                          
   +----------------+             +--------------+             +--------------+             +----------------+  
   |                |     HTTP    |              |    libp2p   |              |     HTTP    |                |  
   |   homeserver   |<----------->| libp2p-proxy |<----------->| libp2p-proxy |<----------->|   homeserver   |  
   |   /ipfs/Qm..   |             |              |             |              |             |   /ipfs/Qm..   |  
   |                |             +--------------+             +--------------+             |                |  
   +----------------+                                                                       +----------------+  
                                                                                                       
```
