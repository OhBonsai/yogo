package api

import (
	"github.com/OhBonsai/yogo/app"
	"github.com/OhBonsai/yogo/model"
	"github.com/OhBonsai/yogo/store/storetest"
	"github.com/OhBonsai/yogo/store/sqlstore"
	"github.com/OhBonsai/yogo/store"
)

type TestHelper struct {
	App            *app.App
	tempConfigPath string

	Client         model.Client
}

type persistentTestStore struct {
	store.Store
}

func (*persistentTestStore) Close() {}

var testStoreContainer *storetest.RunningContainer
var testStore *persistentTestStore

func UseTestStore(container *storetest.RunningContainer, settings *model.SqlSettings) {
	testStoreContainer = container
	testStore = &persistentTestStore{store.NewLayeredStore(sqlstore.NewSqlSupplier(*settings))}
}

func StopTestStore() {
	if testStoreContainer != nil {
		testStoreContainer.Stop()
		testStoreContainer = nil
	}
}

