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
	"log/syslog"
	"math"
	"net/http"
	"os"
	"strconv"
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

	pointdata := loadData()

	for {
		select {
		case <-time.After(time.Second * 1):
			log.Println("Tick:")
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

func loadData() []float32 {
	csvfile, err := os.Open("GDA94_SITE.CSV")

	if err != nil {
		fmt.Println(err)
		return nil
	}

	defer csvfile.Close()
	r := csv.NewReader(csvfile)

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
