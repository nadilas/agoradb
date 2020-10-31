// Copyright 2020. feat.Me Networks. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package schema

import (
	"github.com/jhump/protoreflect/desc"
)

type Database struct {
	Name       string
	Descriptor *desc.FileDescriptor
}
