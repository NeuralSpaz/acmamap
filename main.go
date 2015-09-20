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
	"fmt"
	"log"
	"log/syslog"
	"net/http"

	"golang.org/x/net/websocket"
)

func init() {
	lw, e := syslog.New(syslog.LOG_NOTICE, "ACMAMAP")
	if e == nil {
		log.SetOutput(lw)
	}
}

func main() {
	log.Println("Starting ACMAMAP")
	http.Handle("/", http.FileServer(http.Dir("www/")))
	http.Handle("/ws/acmasites", websocket.Handler(wsACMASites))
	port := 30000
	fmt.Printf("Server running on port %d\n", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		log.Println(err)
	}
}
