package main

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

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"log/syslog"
	"net/http"
	"time"

	"golang.org/x/net/websocket"
)

func init() {
	// Use syslog for all our logging needs
	lw, e := syslog.New(syslog.LOG_NOTICE, "ACMAMAP")
	if e == nil {
		log.SetOutput(lw)
	}
}

func main() {
	// Tell syslog we are starting
	log.Println("Starting ACMAMAP")

	// Register our http Handlers
	http.Handle("/", http.FileServer(http.Dir("www/")))
	// Register out websockets Handlers
	http.Handle("/ws/acmasites", websocket.Handler(wsACMASites))

	// Spin that http server up!
	port := 30000
	fmt.Printf("Server running on port %d\n", port)
	go func() {
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
	}()

	pointdata := []float32{221.04538524444445, 136.93982096462926, 221.0459136, 136.93963306394787, 221.06414648888887, 136.90747494091536, 220.98927644444444, 136.9198768606348, 223.14627413333335, 145.4789690031631}
	// Main Loop
	for {
		select {
		case <-time.After(time.Second * 1):
			log.Println("Tick:")
		case cl := <-newConn:
			buf := new(bytes.Buffer)
			// var pi float64 = math.Pi
			err := binary.Write(buf, binary.LittleEndian, pointdata)
			if err != nil {
				log.Println("binary.Write failed:", err)
			}
			err = websocket.Message.Send(cl.ws, buf.Bytes())
			if err != nil {
				log.Println(err)
			}
			log.Printf("Sent %v bytes", buf.Len())
		}
	}
}
