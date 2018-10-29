package api

import (
	"github.com/OhBonsai/yogo/app"
	"github.com/OhBonsai/yogo/web"
	"github.com/gorilla/mux"
	"github.com/OhBonsai/yogo/model"
	"net/http"
)

type Routes struct {
	Root *mux.Router    // ''
	ApiRoot *mux.Router // '/api/v1'

	Plans *mux.Router   // '/api/v1/plans'
	Plan  *mux.Router   // '/api/v1/plan/{post_id:[A-Za-z0-9]+}'
}

type API struct {
	App        *app.App
	BaseRoutes *Routes
}


func Init(a *app.App, root *mux.Router) *API {
	api := &API{
		App: a,
		BaseRoutes: &Routes{},
	}

	api.BaseRoutes.Root = root
	api.BaseRoutes.ApiRoot = root.PathPrefix(model.API_URL_SUFFIX).Subrouter()

	api.BaseRoutes.Plans = api.BaseRoutes.ApiRoot.PathPrefix("/plans").Subrouter()
	api.BaseRoutes.Plan = api.BaseRoutes.ApiRoot.PathPrefix("/plan/{plan_id:[A-zA-Z0-9]+}").Subrouter()

	api.InitPlan()

	root.Handle("/api/v1/{anything:.*}", http.HandlerFunc(api.Handle404))
	return api
}

func (api *API) Handle404(w http.ResponseWriter, r *http.Request) {
	web.Handle404(api.App, w, r)
}

var ReturnStatusOK = web.ReturnStatusOK