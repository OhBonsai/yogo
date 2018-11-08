package app

import (
	"fmt"
	"net/http"
	"sync/atomic"

	"github.com/OhBonsai/yogo/mlog"
	"github.com/OhBonsai/yogo/store"
	"github.com/OhBonsai/yogo/store/sqlstore"
	"github.com/OhBonsai/yogo/model"
	"github.com/OhBonsai/yogo/utils"
	"github.com/gorilla/mux"
)

type App struct {
	Srv *Server
	Log *mlog.Logger
	newStore func() store.Store

	config          atomic.Value
	envConfig       map[string]interface{}
	configFile      string
}

func New(configFileLocation string) (outApp *App, outErr error) {
	app := &App {
		Srv: &Server{
			Router: mux.NewRouter(),
		},
		configFile: configFileLocation,
	}


	if err := app.LoadConfig(app.configFile); err != nil {
		return nil, err
	}

	// Initalize logging
	app.Log = mlog.NewLogger(utils.MloggerConfigFromLoggerConfig(&app.Config().LogSettings))
	// Redirect default golang logger to this logger
	mlog.RedirectStdLog(app.Log)
	// Use this app logger as the global logger (eventually remove all instances of global logging)
	mlog.InitGlobalLogger(app.Log)


	mlog.Info("Server is initializing...")

	if app.newStore == nil {
		app.newStore = func() store.Store {
			return store.NewLayeredStore(sqlstore.NewSqlSupplier(*new(model.SqlSettings).SetDefaults()))
		}
	}

	app.initJobs()
	app.Srv.Store = app.newStore()
	app.Srv.Router.NotFoundHandler = http.HandlerFunc(app.Handle404)
	return app, nil

}


func (a *App) Handle404(w http.ResponseWriter, r *http.Request) {
	err := model.NewAppError("Handle404", "api.context.404.app_error", nil, "", http.StatusNotFound)
	mlog.Debug(fmt.Sprintf("%v: code=404 ip=%v", r.URL.Path, utils.GetIpAddress(r)))
	fmt.Fprint(w, err.Message)
}


func (a *App) initJobs() {
	// do some background task. like cmdb data sync
}


func (a *App) Shutdown(){

}