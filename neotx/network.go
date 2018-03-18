package neotx

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net"

	"github.com/o3labs/neo-transaction-watcher/neotx/network"
)

//Network
const (
	NEOMainNet    network.NEONetworkMagic = 7630401
	NEOTestNet    network.NEONetworkMagic = 1953787457
	NEOPrivateNet network.NEONetworkMagic = 56753
)

type Config struct {
	Network    network.NEONetworkMagic
	IPAddress  string
	Port       uint16
	connection net.Conn
}

type TX struct {
	Type network.InventoryType
	ID   string
}

type OnReceivedTX func(tx TX)

type Interface interface {
	handleConnection(handler OnReceivedTX)
}

var _ Interface = (*Config)(nil)

func (c *Config) handleConnection(handler OnReceivedTX) {
	conn := c.connection
	log.Printf("remote address = %v", conn.RemoteAddr().String())
	log.Printf("local address = %v", conn.LocalAddr().String())
	nonce, _ := network.RandomUint32()
	payload := network.NewVersionPayload(c.Port, nonce)
	versionCommand := network.NewMessage(c.Network, network.CommandVersion, payload)
	conn.Write(versionCommand)

	for {
		_, msg, err := network.ReadMessage(conn, nil)
		if err != nil {
			log.Printf("%v", err)
			return
		}

		//receive version from remote node
		if msg.Command == string(network.CommandVersion) {
			out := &network.Version{}
			payloadByte := make([]byte, msg.Length)
			_, err = io.ReadFull(conn, payloadByte)
			if err != nil {
				log.Printf("err %v ", err)
				//continue
				continue
			}
			pr := bytes.NewBuffer(payloadByte)
			out.Decode(pr, 0)
			//reply with verack
			verack := network.NewMessage(c.Network, network.CommandVerack, nil)
			conn.Write(verack)

		} else if msg.Command == string(network.CommandVerack) {

		} else if msg.Command == string(network.CommandAddr) {
			out := &network.Addr{}
			payloadByte := make([]byte, msg.Length)
			_, err = io.ReadFull(conn, payloadByte)
			if err != nil {
				log.Printf("err %v ", err)
				continue
			}
			pr := bytes.NewBuffer(payloadByte)
			out.Decode(pr, 0)
		} else if msg.Command == string(network.CommandInv) {
			out := &network.Inv{}
			payloadByte := make([]byte, msg.Length)
			_, err = io.ReadFull(conn, payloadByte)
			if err != nil {
				log.Printf("err %v ", err)
				continue
			}
			pr := bytes.NewBuffer(payloadByte)
			out.Decode(pr, 0)
			b := out.Hash.ToBytes()
			txID := hex.EncodeToString(b)
			tx := TX{
				Type: out.Type,
				ID:   txID,
			}
			go handler(tx)
		}
	}
}

func Start(config Config, handler OnReceivedTX) error {
	if handler == nil {
		return fmt.Errorf("handler cannot be null")
	}
	address := fmt.Sprintf("%v:%v", config.IPAddress, config.Port)
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return err
	}
	config.connection = conn
	config.handleConnection(handler)
	return nil
}
