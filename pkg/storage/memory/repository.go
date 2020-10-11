// Copyright 2020. feat.Me Networks. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package memory

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/buraksezer/olric"
	"github.com/buraksezer/olric/config"
	"github.com/buraksezer/olric/query"
	"github.com/featme-inc/agoradb/internal/services"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/protoparse"
	"github.com/jhump/protoreflect/desc/protoprint"
	"github.com/sirupsen/logrus"
)

var (
	DatabaseAlreadyExists = errors.New("database already exists")
)

type Repository struct {
	db *olric.Olric
	// <string, protodescriptor>
	databases *olric.DMap
}

func New() *Repository {
	repo := &Repository{}

	c := config.New("local")
	c.LogLevel = "WARN"
	c.Logger = log.New(logrus.StandardLogger().Out, "", log.Flags())
	startup, cancel := context.WithCancel(context.Background())
	c.Started = func() {
		defer cancel()
		log.Println("[INFO] Olric is ready to accept connections")
	}

	var err error
	db, err := olric.New(c)
	if err != nil {
		log.Fatalf("Failed to create Olric instance: %v", err)
	}

	go func() {
		err := db.Start()
		if err != nil {
			log.Fatalf("olric.Start returned an error: %v", err)
		}
	}()

	<-startup.Done()

	databases, err := db.NewDMap("databases")
	repo.db = db
	repo.databases = databases

	return repo
}

func (r *Repository) Stop() error {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	log.Print("Schema shutting down...")
	return r.db.Shutdown(ctx)
}

func (r *Repository) AllDatabases() (databases []services.Database) {
	cursor, err := r.databases.Query(query.M{"$onKey": query.M{"$regexMatch": ""}})
	if err != nil {
		log.Fatalf("databases map is not available: %v", err)
	}
	defer cursor.Close()
	cursor.Range(func(databaseKey string, value interface{}) bool {
		accessor := protoparse.FileContentsFromMap(map[string]string{databaseKey: value.(string)})
		p := protoparse.Parser{
			Accessor: accessor,
		}
		files, err := p.ParseFiles(databaseKey)
		if err != nil {
			logrus.Errorf("Failed to decode proto to file descriptor: %v", err)
			return true
		}
		databases = append(databases, services.Database{Name: databaseKey, Descriptor: files[0]})
		return true
	})

	return
}

func (r *Repository) CreateDatabase(database string, fd *desc.FileDescriptor) error {
	p := protoprint.Printer{}
	str, _ := p.PrintProtoToString(fd)
	err := r.databases.PutIfEx(database, str, 0, olric.IfNotFound)
	if err == olric.ErrKeyFound {
		return DatabaseAlreadyExists
	}
	return err
}
