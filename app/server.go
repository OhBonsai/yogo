package app

import (
	"fmt"
	"net/http"
	"net"
	"github.com/gorilla/mux"
	"github.com/OhBonsai/yogo/store"
	"github.com/pkg/errors"
	"github.com/OhBonsai/yogo/mlog"
	"github.com/OhBonsai/yogo/model"
	"github.com/OhBonsai/yogo/utils"
	"strings"
	"time"
	"os"
	"github.com/gorilla/handlers"
)

var allowedMethods = []string{
	"POST",
	"GET",
	"OPTIONS",
	"PUT",
	"PATCH",
	"DELETE",
}

type Server struct {
	Server *http.Server
	Store store.Store
	Router *mux.Router
}

type CorsWrapper struct {
	config model.ConfigFunc
	router *mux.Router
}

type RecoveryLogger struct {
}

func (rl *RecoveryLogger) Println(i ...interface{}) {
	mlog.Error("Please check the std error output for the stack trace")
	mlog.Error(fmt.Sprint(i))
}

func (cw *CorsWrapper) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))

	if r.Method == "OPTIONS" {
		w.Header().Set(
			"Access-Control-Allow-Methods",
			strings.Join(allowedMethods, ", "))

		w.Header().Set(
			"Access-Control-Allow-Headers",
			r.Header.Get("Access-Control-Request-Headers"))
	}

	if r.Method == "OPTIONS" {
		return
	}

	cw.router.ServeHTTP(w, r)
}


func (a *App) StartServer() error{
	mlog.Info("Starting Server...")

	var handler http.Handler = &CorsWrapper{a.Config, a.Srv.Router}
	a.Srv.Server = &http.Server{
		Handler:      handlers.RecoveryHandler(handlers.RecoveryLogger(&RecoveryLogger{}), handlers.PrintRecoveryStack(true))(handler),
		//Handler:      handler,
		ReadTimeout:  time.Duration(300) * time.Second,
		WriteTimeout: time.Duration(300) * time.Second,
		ErrorLog:     a.Log.StdLog(mlog.String("source", "httpserver")),
	}

	listener, err := net.Listen("tcp", *a.Config().ServiceSettings.ListenAddress)
	if err != nil {
		errors.Wrapf(err, utils.T("api.server.start_server.starting.critical"), err)
		return err
	}


	err = a.Srv.Server.Serve(listener)

	if err != nil && err != http.ErrServerClosed {
		mlog.Critical(fmt.Sprintf("Error starting server, err:%v", err))
		os.Exit(1)
	}
	return nil
}
