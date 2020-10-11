// Copyright 2020. feat.Me Networks. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package services

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type ReadDatabaseRepository interface {
	AllDatabases() []Database
}

type Manager struct {
	dbRepository ReadDatabaseRepository
	registry     *registry
}

func NewManager(dbRepository ReadDatabaseRepository) *Manager {
	return &Manager{
		dbRepository: dbRepository,
		registry:     newRegistry(),
	}
}

func (m *Manager) StartDatabases() {
	for _, db := range m.dbRepository.AllDatabases() {
		go m.StartDatabase(db)
	}
}

func (m *Manager) Stop() error {
	logrus.Infof("ServiceManager shutting down...")
	gracefulCh := make(chan interface{})
	serviceMap := m.registry.serviceMap
	allDatabaseCount := len(serviceMap)

	ctx, _ := context.WithTimeout(context.Background(), 30 * time.Second)
	m.registry.Lock()
	defer m.registry.Unlock()
	for _, databaseService := range serviceMap {
		go func() {
			// maybe use parallel passthrough and check sum of all stops
			databaseService.GracefulStop()
			gracefulCh <- 1
		}()
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(10 * time.Millisecond):
			closedDatabaseCount := len(gracefulCh)
			shutdownComplete := closedDatabaseCount == allDatabaseCount
			if shutdownComplete {
				close(gracefulCh)
				return nil
			}
		}
	}

	return nil
}

func (m *Manager) BackendForDatabaseService(fullMethodName string) string {
	logrus.Debugf("Finding database databaseService address for %s", fullMethodName)
	svc := m.registry.serviceByMethod(fullMethodName)
	if svc == nil {
		return ""
	}
	return svc.PassthroughAddress()
}

type databaseInfo struct {
	Name string
	Info map[string]grpc.ServiceInfo
}

func (m Manager) AllDatabases() []databaseInfo {
	var list []databaseInfo
	for _, s := range m.registry.serviceMap {
		list = append(list, databaseInfo{
			Name: s.Name(),
			Info: s.grpcServer.GetServiceInfo(),
		})
	}

	return list
}

func (m *Manager) StartDatabase(database Database) {
	defer func() {
		if err := recover(); err != nil {
			logrus.Errorf("Paniced on database databaseService: %s. Restarting. Error: %v", database.Name, err)
			time.Sleep(50 * time.Millisecond)
			go m.StartDatabase(database)
		}
	}()

	if m.registry.hasService(database.Name) {
		logrus.Errorf("Cannot start already running database databaseService: %s", database.Name)
		return
	}
	svc := newDatabaseService(database)
	m.registry.addService(database.Name, svc)
	if err := svc.ListenAndServe(); err != nil {
		logrus.Errorf("%s.ListenAndServe exited with error: %v", database.Name, err)
	}
	m.registry.removeService(database.Name)
}
