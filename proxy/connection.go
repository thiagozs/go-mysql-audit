package proxy

import (
	"fmt"
	"io"
	"log"
	"net"

	"github.com/thiagozs/go-mysql-audit/protocol"
)

func NewConnection(host string, port string, conn net.Conn,
	id uint64, enableDecoding, verbose bool) *Connection {
	return &Connection{
		host:           host,
		port:           port,
		conn:           conn,
		id:             id,
		enableDecoding: enableDecoding,
		verbose:        verbose,
	}
}

type Connection struct {
	id             uint64
	conn           net.Conn
	host           string
	port           string
	enableDecoding bool
	verbose        bool
}

func (r *Connection) Handle() error {
	address := fmt.Sprintf("%s%s", r.host, r.port)
	mysql, err := net.Dial("tcp", address)
	if err != nil {
		log.Printf("Failed to connection to MySQL: [%d] %s", r.id, err.Error())
		return err
	}

	if !r.enableDecoding {
		// client to server
		go func() {
			copied, err := io.Copy(mysql, r.conn)
			if err != nil {
				log.Printf("Conection error: [%d] %s", r.id, err.Error())
			}

			log.Printf("Connection closed. Bytes copied: [%d] %d", r.id, copied)
		}()

		copied, err := io.Copy(r.conn, mysql)
		if err != nil {
			log.Printf("Connection error: [%d] %s", r.id, err.Error())
		}

		log.Printf("Connection closed. Bytes copied: [%d] %d", r.id, copied)

		return nil
	}

	handshakePacket := &protocol.InitialHandshakePacket{}
	err = handshakePacket.Decode(mysql)
	if err != nil {
		log.Printf("Failed ot decode handshake initial packet: %s", err.Error())
		return err
	}

	//fmt.Printf("InitialHandshakePacket:\n%s\n", handshakePacket)

	res, _ := handshakePacket.Encode()

	written, err := r.conn.Write(res)
	if err != nil {
		log.Printf("Failed to write %d: %s", written, err.Error())
		return err
	}

	go func() {
		copied, err := io.Copy(mysql, r.conn)
		if err != nil {
			log.Printf("Conection error: [%d] %s", r.id, err.Error())
		}

		log.Printf("Connection closed. Bytes copied: [%d] %d", r.id, copied)
	}()

	var copied int64
	if r.verbose {
		if err := proxyLog(r.conn, mysql, r.verbose); err != nil {
			log.Printf("Connection error: [%d] %s", r.id, err.Error())
		}
	} else {
		copied, err = io.Copy(r.conn, mysql)
		if err != nil {
			log.Printf("Connection error: [%d] %s", r.id, err.Error())
		}
	}

	log.Printf("Connection closed. Bytes copied: [%d] %d", r.id, copied)

	return nil
}
