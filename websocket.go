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
	"log"

	"golang.org/x/net/websocket"
)

// Something to hold our clients in
type ClientConn struct {
	ws *websocket.Conn
	// client command chan
	cmd chan uint8
}

// All the Clients
var newConn = make(chan *ClientConn)

func wsACMASites(ws *websocket.Conn) {
	c := &ClientConn{}
	c.ws = ws
	c.cmd = make(chan uint8)
	var cmd uint8
	log.Println("New wsACMASites Conection: ", ws.RemoteAddr())

	newConn <- c

	for {
		// do somthing here
		c.cmd <- cmd
	}

}
