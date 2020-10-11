// Copyright 2020. feat.Me Networks. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package agoradb

// The connection
type Client struct {
	config Config
}

func New(config Config) (*Client, error) {
	cfg, err := mergeConfig(config)
	if err != nil {
		return nil, err
	}
	return &Client{
		config: cfg,
	}, nil
}

// Returns a copy of the internal configuration
func (c	*Client) Config() Config {
	return c.config
}