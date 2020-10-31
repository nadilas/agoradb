// Copyright 2020. feat.Me Networks. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package migration

type Migrate interface {
	Migrate() error
	Revert() error
	Commit() error
}

type migrate struct {
}

func (m *migrate) Migrate() error {
	panic("implement me")
}

func (m *migrate) Revert() error {
	panic("implement me")
}

func (m *migrate) Commit() error {
	panic("implement me")
}
