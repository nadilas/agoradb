// Copyright 2020. feat.Me Networks. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package agoradb


type Auth struct {

	// User to be used for authentication
	User string

	// The database Password for the named user
	Password string

	// The client certificate to be used for authentication
	//
	// If the CertificatePath is set, user and password entries will be ignored.
	CertificatePath string

	// Optionally you can set the RootCA cert, it overrides the system certificates
	RootCA string
}

type Config struct {
	// The agoradb node endpoint to connect to
	//
	// The format is host:port
	Uri    string

	// The authentication to be used
	Auth   Auth
}

var (
	defaultCfg = Config{
		Uri: "localhost:5750",
		Auth: Auth{
			User: "guest",
			Password: "guest",
		},
	}
)

func mergeConfig(config Config) (Config, error) {
	cfg := defaultCfg
	if config.Uri != "" {
		// TODO validate uri scheme

		cfg.Uri = config.Uri
	}
	return cfg, nil
}

