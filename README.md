# neo-transaction-watcher
A headless daemon to connect to the NEO network to receive TX message.

Change the endpoint of the network that you want to connect to and run. You will be able to see transaction coming in in near real-time.

try `go run main.go` 

#### main.go
```go
package main

import (
	"log"
	"os"

	"github.com/o3labs/neo-transaction-watcher/neotx"
	"github.com/o3labs/neo-transaction-watcher/neotx/network"
)

type Handler struct {
}

//implement the message protocol
func (h *Handler) OnReceive(tx neotx.TX) {
	log.Printf("%+v", tx)
}

func (h *Handler) OnConnected(c network.Version) {
	log.Printf("connected %+v", c)
}

func (h *Handler) OnError(e error) {
	log.Printf("error %+v", e)
}

func main() {
	config := neotx.Config{
		Network:   neotx.NEOMainNet,
		Port:      10333,
		IPAddress: "52.193.202.2",
	}
	client := neotx.NewClient(config)
	handler := &Handler{}
	client.SetDelegate(handler)

	err := client.Start()
	if err != nil {
		log.Printf("%v", err)
		os.Exit(-1)
	}

	for {

	}
}

```


You can fetch transaction detail by calling `getrawtransaction` JSON-RPC method.  
You can use [neo-utils](https://github.com/O3Labs/neo-utils) to parse the invocation script to get information about the transaction like script hash, smart contract's method and params.


#### Things that you can build when you have this  
- a notification service to notify user by email about the transaction
- watch particular smart contract
- watch particular NEO address
- etc.
