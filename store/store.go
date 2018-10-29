package store

import (
	"github.com/OhBonsai/yogo/model"
)


type StoreResult struct {
	Data interface{}
	Err *model.AppError
}


type StoreChannel chan StoreResult


type Store interface {
	Close()
	DropAllTables()
	TotalMasterDbConnections() int
	TotalReadDbConnections() int
	TotalSearchDbConnections() int


	Plan() PlanStore
}


type PlanStore interface {
	Save(plan *model.Plan) StoreChannel
	Update(plan *model.Plan) StoreChannel
	Get(id string) StoreChannel
	Delete(planId string, time int64, deleteByID string) StoreChannel
}


func Do(f func(result *StoreResult)) StoreChannel {
	storeChannel := make(StoreChannel, 1)
	go func() {
		result := StoreResult{}
		f(&result)
		storeChannel <- result
		close(storeChannel)
	}()
	return storeChannel
}
