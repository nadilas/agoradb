// Copyright 2020. feat.Me Networks. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package services

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
	"sync"
)

type registry struct {
	mux        sync.RWMutex
	services   [][]byte
	serviceMap map[int]*databaseService
}

func newRegistry() *registry {
	return &registry{
		mux:        sync.RWMutex{},
		serviceMap: make(map[int]*databaseService),
	}
}

func (r *registry) addService(serviceName string, svc *databaseService) int {
	r.mux.Lock()
	defer r.mux.Unlock()
	idx := len(r.services)
	r.services = append(r.services, []byte(serviceName))
	r.serviceMap[idx] = svc
	return idx
}

func (r *registry) removeService(serviceName string) {
	r.mux.Lock()
	defer r.mux.Unlock()
	idx := r.idxOfService(serviceName)
	oldSize := len(r.services)
	r.services = append(r.services[0:idx], r.services[idx+1:]...)
	// shift map
	for i := idx + 1; i < oldSize; i++ {
		r.serviceMap[i-1] = r.serviceMap[i]
	}
	delete(r.serviceMap, oldSize-1)
}

func (r *registry) hasService(serviceName string) bool {
	return r.idxOfService(serviceName) >= 0
}

func (r *registry) idxOfService(serviceName string) int {
	r.mux.RLock()
	defer r.mux.RUnlock()
	for idx, service := range r.services {
		if bytes.ContainsAny(service, serviceName) {
			return idx
		}
	}
	return -1
}

// Find a databaseService instance by calling the FQDN of a method, e.g.: /testdb.User/Save
func (r *registry) serviceByMethod(fullMethodName string) *databaseService {
	dbName := strings.Split(strings.ReplaceAll(fullMethodName, "/", ""), ".")[0]
	rxp := regexp.MustCompile(fmt.Sprintf("^%s$", dbName))
	r.mux.RLock()
	defer r.mux.RUnlock()
	for idx, svcName := range r.services {
		if rxp.Match(svcName) {
			return r.serviceMap[idx]
		}
	}
	return nil
}

func (r *registry) Lock() {
	r.mux.Lock()
}

func (r *registry) Unlock() {
	r.mux.Unlock()
}
