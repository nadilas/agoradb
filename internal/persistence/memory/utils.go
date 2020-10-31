// Copyright 2020. feat.Me Networks. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package memory

import (
	"github.com/buraksezer/olric/query"
)

func keyExistsHelper(key string) query.M {
	return query.M{"$onKey": query.M{"$regexMatch": key, "$options": query.M{"$onValue": query.M{"$ignore": true}}}}
}
