// Copyright 2020. feat.Me Networks. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"log"

	"github.com/featme-inc/agoradb"
	"github.com/featme-inc/agoradb/examples/basic/migrations"
	"github.com/featme-inc/agoradb/internal/schema/migration"
)

func migrate(client *agoradb.Client) {
	m := migration.New("initial")
	err := m.Apply(client,
		&migrations.CreateDatabase{},
		&migrations.ChangeUserProps{},
	)
	if err != nil {
		err = m.Revert()
		if err != nil {
			log.Fatal(err)
		} else {
			log.Print("Failed to migrate, changes have been reverted")
		}
	} else {
		log.Printf("Successfully applied: %s", m.Name())
	}
}
