package web

import (
	"strings"
	"net/http"
	"github.com/OhBonsai/yogo/model"
	"github.com/OhBonsai/yogo/mlog"
	"github.com/OhBonsai/yogo/utils"
	"github.com/OhBonsai/yogo/app"
	"fmt"
	"github.com/gorilla/mux"
)


type Web struct {
	App        *app.App
	MainRouter *mux.Router
}

func Handle404(a *app.App, w http.ResponseWriter, r *http.Request) {
	err := model.NewAppError("Handle404", "api.context.404.app_error", nil, "", http.StatusNotFound)

	mlog.Debug(fmt.Sprintf("%v: code=404 ip=%v", r.URL.Path, utils.GetIpAddress(r)))

	if IsApiCall(r) {
		w.WriteHeader(err.StatusCode)
		err.DetailedError = "There doesn't appear to be an api call for the url='" + r.URL.Path + "'.  Typo? are you missing a team_id or user_id as part of the url?"
		w.Write([]byte(err.ToJson()))
	}else{
		w.WriteHeader(err.StatusCode)
		err.DetailedError = "This is not a api call for the url='" + r.URL.Path + "'. It should start with /api"
		w.Write([]byte(err.DetailedError))
	}
}

func IsApiCall(r *http.Request) bool {
	return strings.Index(r.URL.Path, "/api/") == 0
}


func ReturnStatusOK(w http.ResponseWriter) {
	m := make(map[string]string)
	m[model.STATUS] = model.STATUS_OK
	w.Write([]byte(model.MapToJson(m)))
}
