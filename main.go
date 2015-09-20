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
	"log"
	"log/syslog"
)

func init() {
	logwriter, e := syslog.New(syslog.LOG_NOTICE, "ACMAMAP")
	if e == nil {
		log.SetOutput(logwriter)
	}
}

func main() {
	log.Println("Starting ACMAMAP")
}
