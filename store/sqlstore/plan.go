package sqlstore

import (
	"net/http"

	"github.com/OhBonsai/yogo/model"
	"github.com/OhBonsai/yogo/utils"
	"github.com/OhBonsai/yogo/store"
	"sync"
)


const (
	LAST_PLAN_TIME_CACHE_SIZE = 25000
	LAST_PLAN_TIME_CACHE_SEC  = 900 // 15 minutes

	LAST_PLANS_CACHE_SIZE = 1000
	LAST_PLANS_CACHE_SEC  = 900 // 15 minutes
)



type SqlPlanStore struct {
	SqlStore
	lastPlanTimeCache *utils.Cache
	lastPlansCache    *utils.Cache
	maxPlanSizeOnce   sync.Once
	maxPlanSizeCached int
}


func (s *SqlPlanStore) ClearCaches() {
	s.lastPlanTimeCache.Purge()
	s.lastPlansCache.Purge()
}


func NewSqlPlanStore(sqlStore SqlStore) store.PlanStore {
	s := &SqlPlanStore{
		SqlStore:          sqlStore,
		lastPlanTimeCache: utils.NewLru(LAST_PLAN_TIME_CACHE_SIZE),
		lastPlansCache:    utils.NewLru(LAST_PLANS_CACHE_SIZE),
		maxPlanSizeCached: model.PLAN_MESSAGE_MAX_RUNES_V1,
	}

	for _, db := range sqlStore.GetAllConns() {
		// 这里通过反射会拿到所有字段
		table := db.AddTableWithName(model.Plan{}, "Plans").SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(26)
		table.ColMap("Name").SetMaxSize(model.PLAN_NAME_MAX_BYTES_V1)
		table.ColMap("Props").SetMaxSize(8000)
	}

	return s
}


func (s *SqlPlanStore) CreateIndexesIfNotExists() {
	s.CreateIndexIfNotExists("idx_plans_update_at", "Plans", "UpdateAt")
	s.CreateIndexIfNotExists("idx_plans_create_at", "Plans", "CreateAt")
	s.CreateIndexIfNotExists("idx_plans_delete_at", "Plans", "DeleteAt")
	s.CreateFullTextIndexIfNotExists("idx_plans_name_txt", "Plans", "Name")
}


func (s *SqlPlanStore) Save(plan *model.Plan) store.StoreChannel {
	return store.Do(func(result *store.StoreResult) {
		plan.PreSave()

		if result.Err = plan.IsValid(); result.Err != nil {
			return
		}

		if err := s.GetMaster().Insert(plan); err != nil {
			result.Err = model.NewAppError("SqlPlanStore.Save", "store.sql_plan.save.app_error", nil, "id="+plan.Id+", "+err.Error(), http.StatusInternalServerError)
		}

		result.Data = plan
	})
}



func (s *SqlPlanStore) Update(plan *model.Plan) store.StoreChannel {
	return store.Do(func(result *store.StoreResult) {
		plan.UpdateAt = model.GetMillis()
		plan.PreCommit()

		if result.Err = plan.IsValid(); result.Err != nil {
			return
		}

		if _, err := s.GetMaster().Update(plan); err != nil {
			result.Err = model.NewAppError("SqlPostStore.Update", "store.sql_post.update.app_error", nil, "id="+plan.Id+", "+err.Error(), http.StatusInternalServerError)
		} else {
			result.Data = plan
		}
	})
}


func (s *SqlPlanStore) Get(id string) store.StoreChannel {
	return store.Do(func(result *store.StoreResult) {
		if len(id) == 0 {
			result.Err = model.NewAppError("SqlPlanStore.GetPlan", "store.sql_plan.get.app_error", nil, "id="+id, http.StatusBadRequest)
			return
		}

		var plan model.Plan
		if err := s.GetReplica().SelectOne(&plan, "SELECT * FROM plans WHERE Id = :Id AND DeleteAt = 0", map[string]interface{}{"Id": id}); err != nil {
			result.Err = model.NewAppError("SqlPlanStore.GetPlan", "store.sql_plan.get.app_error", nil, "id="+id+err.Error(), http.StatusNotFound)
		}

		result.Data = &plan
	})
}

func (s *SqlPlanStore) Delete(planId string, time int64, deleteByID string) store.StoreChannel {
	return store.Do(func(result *store.StoreResult) {

		appErr := func(errMsg string) *model.AppError {
			return model.NewAppError("SqlPlanStore.Delete", "store.sql_plan.delete.app_error", nil, "id="+planId+", err="+errMsg, http.StatusInternalServerError)
		}


		var plan model.Plan
		if err := s.GetReplica().SelectOne(&plan, "SELECT * FROM plans WHERE Id = :Id AND DeleteAt = 0", map[string]interface{}{"Id": planId}); err != nil {
			result.Err = appErr(err.Error())
		}


		plan.AddProp(model.PLAN_PROPS_DELETE_BY, deleteByID)
		if _, err := s.GetMaster().Exec("UPDATE Plans SET DeleteAt = :DeleteAt, UpdateAt = :UpdateAt, Props = :Props WHERE Id = :Id", map[string]interface{}{"DeleteAt": time, "UpdateAt": time, "Id": planId, "Props": model.StringInterfaceToJson(plan.Props)}); err != nil {
			result.Err = appErr(err.Error())
		}
	})
}


