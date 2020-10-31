// Copyright 2020. feat.Me Networks. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package schema

import (
	"context"

	"github.com/featme-inc/agoradb/internal/schema/schemapb"
)

func (s *Server) DescribeDatabase(ctx context.Context, database *schemapb.Database) (*schemapb.DescribeDatabaseResponse, error) {
	panic("implement me")
}
