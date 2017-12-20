package main

import (
	"log"
	"net"
)

const protocol = "tcp"

var fullNodeAddress string
var knownNodes = []string{"172.19.0.2:8081", "172.19.0.3:8082", "172.19.0.4:8083", "172.19.0.5:8084"} //IP and ports of our peers/nodes

//StartServer should connect to, and start a node
func StartServer(peerAddress string) {
	fullNodeAddress = peerAddress
	_, port, _ := net.SplitHostPort(peerAddress)

	listen, error := net.Listen(protocol, fullNodeAddress)
	if error != nil {
		panic(error)
	}
	defer listen.Close()

	bc := NewBlockChain(port)

	for {
		connection, err := listen.Accept()
		if err != nil {
			log.Panic(err)
		}
		go handleConnection(connection, bc)

	}

}

func handleConnection(conn net.Conn, bc *Blockchain) {

}
