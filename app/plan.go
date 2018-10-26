package app

import (
	"healer/model"
)


func (a *App) CreatePlan(plan *model.Plan) (*model.Plan, *model.AppError) {
	if result := <-a.Srv.Store.Plan().Save(plan); result.Err != nil {
		return nil, result.Err
	} else {
		rplan := result.Data.(*model.Plan)
		return rplan, nil
	}
}


func (a *App) UpdatePlan(plan *model.Plan) (*model.Plan, *model.AppError) {

	if result := <-a.Srv.Store.Plan().Get(plan.Id); result.Err != nil {
		return nil, result.Err
	}

	if result := <-a.Srv.Store.Plan().Update(plan); result.Err != nil {
		return nil, result.Err
	} else {
		rplan := result.Data.(*model.Plan)
		return rplan, nil
	}

}