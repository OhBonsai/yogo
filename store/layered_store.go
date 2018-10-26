package store

import (
	"context"
)

type LayeredStoreDatabaseLayer interface {
	LayeredStoreSupplier
	Store
}

type LayeredStore struct {
	TmpContext      context.Context
	DatabaseLayer   LayeredStoreDatabaseLayer
}



func NewLayeredStore(db LayeredStoreDatabaseLayer) Store {
	ls := &LayeredStore{
		TmpContext:      context.TODO(),
		DatabaseLayer:   db,
	}

	return ls
}


func (s *LayeredStore) Close() {
	s.DatabaseLayer.Close()
}

func (s *LayeredStore) DropAllTables() {
	s.DatabaseLayer.DropAllTables()
}

func (s *LayeredStore) TotalMasterDbConnections() int {
	return s.DatabaseLayer.TotalMasterDbConnections()
}


func (s *LayeredStore) TotalReadDbConnections() int {
	return s.DatabaseLayer.TotalReadDbConnections()
}

func (s *LayeredStore) TotalSearchDbConnections() int {
	return s.DatabaseLayer.TotalSearchDbConnections()
}

func (s *LayeredStore) Plan() PlanStore {
	return s.DatabaseLayer.Plan()
}


