package main

import (
	"github.com/julianknutsen/wasteland/internal/backend"
	"github.com/julianknutsen/wasteland/internal/commons"
)

// openStore creates a WLCommonsStore for the given local database directory.
// Package-level variable to allow test overrides.
var openStore = func(localDir string, signed bool, hopURI string) commons.WLCommonsStore {
	db := backend.NewLocalDB(localDir, "")
	store := commons.NewWLCommons(db)
	store.SetSigning(signed)
	store.SetHopURI(hopURI)
	return store
}

// openDB creates a commons.DB for the given local database directory.
// Package-level variable to allow test overrides.
var openDB = func(localDir string) commons.DB {
	return backend.NewLocalDB(localDir, "")
}
