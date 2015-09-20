// ACMA Data Mapping
//  Copyright (C) 2015  Josh Gardiner aka github.com/NeuralSpaz/

//  This program is free software; you can redistribute it and/or modify
//  it under the terms of the GNU General Public License as published by
//  the Free Software Foundation; either version 2 of the License, or
//  (at your option) any later version.

//  This program is distributed in the hope that it will be useful,
//  but WITHOUT ANY WARRANTY; without even the implied warranty of
//  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
//  GNU General Public License for more details.

//  You should have received a copy of the GNU General Public License along
//  with this program; if not, write to the Free Software Foundation, Inc.,
//  51 Franklin Street, Fifth Floor, Boston, MA 02110-1301 USA.package main

// Handler for our binary websockets data connection

package main

import (
	"bytes"
	"encoding/binary"
	"log"
	"sync"

	"golang.org/x/net/websocket"
)

// Something to hold our clients in
type ClientConn struct {
	ws       *websocket.Conn
	inBuffer [100]byte
	cmd      chan uint8
}

// All the Clients
var newConn = make(chan *ClientConn)

type counter struct {
	sync.Mutex
	Clients int
}

var clients counter

func wsACMASites(ws *websocket.Conn) {
	c := &ClientConn{}
	c.ws = ws
	c.cmd = make(chan uint8)
	var cmd uint8
	log.Println("New wsACMASites Conection: ", ws.RemoteAddr())

	newConn <- c
	clients.Lock()
	clients.Clients++
	clients.Unlock()

	for {
		// read the packet
		// TODO: Reduce allocations
		packet := c.inBuffer[0:]
		n, err := ws.Read(packet)
		packet = packet[0:n]
		if err != nil {
			log.Println("wsACMASites Read Error: ", err)
			clients.Lock()
			clients.Clients--
			clients.Unlock()
			break
		}
		buf := bytes.NewBuffer(packet)
		// its all binary
		err = binary.Read(buf, binary.LittleEndian, &cmd)
		if err != nil {
			log.Println(err)
			clients.Lock()
			clients.Clients--
			clients.Unlock()
			break
		}
		// Hope this works like I think it will
		c.cmd <- cmd
	}

}
