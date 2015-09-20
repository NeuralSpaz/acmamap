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
	"encoding/csv"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
	"time"

	"golang.org/x/net/websocket"
)

// func init() {
// 	// Use syslog for all our logging needs
// 	lw, e := syslog.New(syslog.LOG_NOTICE, "ACMAMAP")
// 	if e == nil {
// 		log.SetOutput(lw)
// 	}
// }

func main() {
	// Tell syslog we are starting
	log.Println("Starting ACMAMAP")

	// Register our http Handlers
	http.HandleFunc("/", handleIndex)
	// http.Handle("/", http.FileServer(http.Dir("www/")))
	// Register out websockets Handlers
	http.Handle("/ws/acmasites", websocket.Handler(wsACMASites))

	// Spin that http server up!
	port := 30000
	fmt.Printf("Server running on port %d\n", port)
	go func() {
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
	}()

	pointdata := loadData()
	log.Println("Loaded Point Data")
	for {
		select {
		case <-time.After(time.Second * 10):
			clients.Lock()
			log.Println("Current Clients: ", clients.Clients)
			clients.Unlock()
		case cl := <-newConn:
			buf := new(bytes.Buffer)
			// var pi float64 = math.Pi
			err := binary.Write(buf, binary.LittleEndian, pointdata[2:])
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

func handleIndex(rw http.ResponseWriter, req *http.Request) {
	if req.Method != "GET" {
		http.Error(rw, "opps", http.StatusBadRequest)
		log.Printf("%+v\n", req)
	}

	log.Println("Serving JavaScript")
	request := req.URL.Path[1:]
	log.Println(request)

	if request == "" {
		request = "index.html"
	}
	file, err := Asset("www/" + request)
	if err != nil {
		log.Printf("%+v\n", req)
		log.Println(err)
		http.NotFound(rw, req)
	}
	rw.Write(file)
}

func loadData() []float32 {

	csvfile, err := Asset("www/GDA94_SITE.CSV")
	if err != nil {
		fmt.Println(err)
		return nil
	}
	file := bytes.NewReader(csvfile)

	r := csv.NewReader(file)

	records, err := r.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	var points []float32

	for _, each := range records {
		lat, _ := strconv.ParseFloat(each[1], 64)
		lng, _ := strconv.ParseFloat(each[2], 64)
		lng = 256 * (0.5 + lng/360)
		lat = project(lat)
		lat32, lng32 := float32(lat), float32(lng)
		points = append(points, lng32)
		points = append(points, lat32)
	}
	return points
}

func project(lat float64) float64 {
	siny := math.Sin(lat * math.Pi / 180)
	siny = math.Min(math.Max(siny, -0.9999), 0.9999)
	return (256 * (0.5 - math.Log((1+siny)/(1-siny))/(4*math.Pi)))
}
