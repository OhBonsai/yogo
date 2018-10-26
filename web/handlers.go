package web

import (
	"net/http"
	"fmt"
	"healer/app"
	"healer/mlog"
	"healer/model"
	"healer/utils"
)

type Handler struct {
	App            *app.App
	HandleFunc     func(*Context, http.ResponseWriter, *http.Request)
	RequireSession bool
	TrustRequester bool
	RequireMfa     bool
	IsStatic       bool
}


func (w *Web) NewHandler(h func(*Context, http.ResponseWriter, *http.Request)) http.Handler {
	return &Handler{
		App:            w.App,
		HandleFunc:     h,
		RequireSession: false,
		TrustRequester: false,
		RequireMfa:     false,
		IsStatic:       false,
	}
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mlog.Debug(fmt.Sprintf("%v - %v", r.Method, r.URL.Path))

	c := &Context{}
	c.App = h.App
	c.T = utils.T
	c.RequestId = model.NewId()
	c.IpAddress = utils.GetIpAddress(r)
	c.Params = ParamsFromRequest(r)
	c.Path = r.URL.Path
	c.Log = c.App.Log

	// All api response bodies will be JSON formatted by default
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "GET" {
		w.Header().Set("Expires", "0")
	}


	c.Log = c.App.Log.With(
		mlog.String("path", c.Path),
		mlog.String("request_id", c.RequestId),
		mlog.String("ip_addr", c.IpAddress),
		mlog.String("user_id", c.Session.UserId),
		mlog.String("method", r.Method),
	)


	if c.Err == nil {
		h.HandleFunc(c, w, r)
	}

	// Handle errors that have occurred
	if c.Err != nil {
		c.Err.Translate(c.T)
		c.Err.RequestId = c.RequestId

		if c.Err.Id == "api.context.session_expired.app_error" {
			c.LogInfo(c.Err)
		} else {
			c.LogError(c.Err)
		}

		c.Err.Where = r.URL.Path
		w.WriteHeader(c.Err.StatusCode)
		w.Write([]byte(c.Err.ToJson()))

	}

}