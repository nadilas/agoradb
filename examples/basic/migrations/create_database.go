// Copyright 2020. feat.Me Networks. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package migrations

type CreateDatabase struct {}

func (t *CreateDatabase) Migrate() error {
	panic("implement me")
}

func (t *CreateDatabase) Revert() error {
	panic("implement me")
}

func (t *CreateDatabase) Commit() error {
	panic("implement me")
}