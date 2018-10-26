package api

import (
	"healer/model"
	"net/http"
)


func (api *API) InitPlan() {
	api.BaseRoutes.Plans.Handle("", api.ApiHandler(createPlan)).Methods("POST")
	api.BaseRoutes.Plans.Handle("", api.ApiHandler(getPlans)).Methods("GET")

	api.BaseRoutes.Plan.Handle("", api.ApiHandler(getPlan)).Methods("GET")
	api.BaseRoutes.Plan.Handle("", api.ApiHandler(updatePlan)).Methods("PUT")
	api.BaseRoutes.Plan.Handle("", api.ApiHandler(deletePlan)).Methods("DELETE")
}


func createPlan(c *Context, w http.ResponseWriter, r *http.Request) {
	plan := model.PlanFromJson(r.Body)

	if plan == nil {
		c.SetInvalidParam("plan")
		return
	}

	rp, err := c.App.CreatePlan(plan)
	if err != nil {
		c.Err = err
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(rp.ToJson()))
}


func getPlans(c *Context, w http.ResponseWriter, r *http.Request) {

}


func getPlan(c *Context, w http.ResponseWriter, r *http.Request) {

}


func updatePlan(c *Context, w http.ResponseWriter, r *http.Request) {

}


func deletePlan(c *Context, w http.ResponseWriter, r *http.Request) {

}