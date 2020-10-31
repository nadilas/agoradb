// Copyright 2020. feat.Me Networks. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package migration

import "github.com/featme-inc/agoradb"

type Migration interface {
	// Should return the summary of the migration
	Name() string
	Apply(conn *agoradb.Client, migrations ...Migrate) error
	Revert(migrations ...Migrate) error
}

func New(change string) Migration {
	return &migrator{
		name: change,
	}
}

type migrator struct {
	name string
}

func (m *migrator) Name() string {
	return m.name
}

func (m *migrator) Apply(conn *agoradb.Client, migrations ...Migrate) error {
	panic("implement me")
}

func (m *migrator) Revert(migrations ...Migrate) error {
	panic("implement me")
}
