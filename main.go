package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net"

	"github.com/o3labs/neod/network"
)

const (
	port = 30333
)

func handleConnection(conn net.Conn) {
	log.Printf("remote address = %v", conn.RemoteAddr().String())
	log.Printf("local address = %v", conn.LocalAddr().String())
	nonce, _ := network.RandomUint32()
	payload := network.NewVersionPayload(port, nonce)
	versionCommand := network.NewMessage(network.NEOMagic, network.CommandVersion, payload)
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
			}
			pr := bytes.NewBuffer(payloadByte)
			out.Decode(pr, 0)
			//reply with verack
			verack := network.NewMessage(network.NEOMagic, network.CommandVerack, nil)
			conn.Write(verack)

		} else if msg.Command == string(network.CommandVerack) {

		} else if msg.Command == string(network.CommandAddr) {
			out := &network.Addr{}
			payloadByte := make([]byte, msg.Length)
			_, err = io.ReadFull(conn, payloadByte)
			if err != nil {
				log.Printf("err %v ", err)
				//continue
			}
			pr := bytes.NewBuffer(payloadByte)
			out.Decode(pr, 0)
		} else if msg.Command == string(network.CommandInv) {
			out := &network.Inv{}
			payloadByte := make([]byte, msg.Length)
			_, err = io.ReadFull(conn, payloadByte)
			if err != nil {
				log.Printf("err %v ", err)
				//continue
			}
			pr := bytes.NewBuffer(payloadByte)
			out.Decode(pr, 0)
			b := out.Hash.ToBytes()
			if out.Type == network.InventotyTypeTX {
				//You can fetch getrawtransaction to see the transaction detail by txid
				fmt.Printf("type = %v \ninv %v\n", out.Type, hex.EncodeToString(b))
			}

		}

	}
}

func startServer() {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		fmt.Println(err)
		return
	}
	conn, err := ln.Accept()
	if err != nil {
		fmt.Println(err)
		return
	}
	//if remote node is trying to connect to this port then reply with local node version
	log.Printf("accepted\n")
	//once accept connection we then declear version.
	nonce, _ := network.RandomUint32()
	payload := network.NewVersionPayload(port, nonce)
	versionCommand := network.NewMessage(network.NEOMagic, network.CommandVersion, payload)
	conn.Write(versionCommand)

	handleConnection(conn)
}

func startConnectToSeed() {
	// connect to this socket
	conn, err := net.Dial("tcp", "localhost:30333")
	if err != nil {
		fmt.Println(err)
		return
	}
	handleConnection(conn)
}

func main() {

	// go startServer()
	go startConnectToSeed()
	for {

	}
}
