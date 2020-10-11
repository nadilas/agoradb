// Copyright 2020. feat.Me Networks. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"log"

	"github.com/featme-inc/agoradb"
)

// The following go generate creates the client
//go:generate agoradb generate client --addr=loclhost:5750 --auth=BearerToken --grpc --output_dir=basic basic

// This program demonstrates the basic functionality of agoradb
//
// Main focus is:
// 1. client setup
// 2. migration
// 3.
func main() {
	// connect to the database
	username, password := "", ""
	fmt.Print("Username:")
	_, err := fmt.Scan(&username)
	fatalOnError(err)
	fmt.Print("Password:")
	_, err = fmt.Scan(&username)
	fatalOnError(err)

	// Connect with user & pass
	client, err := agoradb.New(agoradb.Config{
		Auth: agoradb.Auth{
			User: username,
			Password: password,
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	// run migrations
	migrate(client)

	// run demo actions
	demo(client)
}
