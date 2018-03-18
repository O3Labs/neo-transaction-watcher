package main

import (
	"log"
	"os"

	"github.com/o3labs/neo-transaction-watcher/neotx"
)

//this method conforms the interface
func OnReceived(tx neotx.TX) {
	log.Printf("%+v", tx)
}

func main() {
	config := neotx.Config{
		Network:   neotx.NEOMainNet,
		Port:      10333,
		IPAddress: "52.193.202.2",
	}
	err := neotx.Start(config, OnReceived)
	if err != nil {
		log.Printf("%v", err)
		os.Exit(-1)
	}

	for {

	}
}
