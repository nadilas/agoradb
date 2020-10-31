// Copyright 2020. feat.Me Networks. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package migration

import (
	"reflect"
	"testing"
)

func TestNewMigration(t *testing.T) {
	tests := []struct {
		name string
		want Migration
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New("dummy_change"); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_migrate_Commit(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &migrate{}
			if err := m.Commit(); (err != nil) != tt.wantErr {
				t.Errorf("Commit() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_migrate_Migrate(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &migrate{}
			if err := m.Migrate(); (err != nil) != tt.wantErr {
				t.Errorf("Migrate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_migrate_Revert(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &migrate{}
			if err := m.Revert(); (err != nil) != tt.wantErr {
				t.Errorf("Revert() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_migrator_Apply(t *testing.T) {
	type args struct {
		migrations []Migrate
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &migrator{}
			if err := m.Apply(tt.args.migrations, nil); (err != nil) != tt.wantErr {
				t.Errorf("Apply() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_migrator_Revert(t *testing.T) {
	type args struct {
		migrations []Migration
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &migrator{}
			if err := m.Revert(tt.args.migrations); (err != nil) != tt.wantErr {
				t.Errorf("Revert() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
