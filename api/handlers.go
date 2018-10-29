package api

import (
	"github.com/OhBonsai/yogo/web"
	"net/http"
)

type Context = web.Context


func (api *API) ApiHandler(h func(*Context, http.ResponseWriter, *http.Request)) http.Handler {
	return &web.Handler{
		App:            api.App,
		HandleFunc:     h,
		RequireSession: false,
		TrustRequester: false,
		RequireMfa:     false,
		IsStatic:       false,
	}
}