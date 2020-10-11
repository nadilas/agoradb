// Copyright 2020. feat.Me Networks. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import "log"

func fatalOnError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
func warnOnError(err error) {
	if err != nil {
		log.Printf("WARN: %v", err)
	}
}
