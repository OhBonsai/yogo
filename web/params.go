package web

import (
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

const (
	PAGE_DEFAULT          = 0
	PER_PAGE_DEFAULT      = 60
	PER_PAGE_MAXIMUM      = 200
	LOGS_PER_PAGE_DEFAULT = 10000
	LOGS_PER_PAGE_MAXIMUM = 10000
)

type Params struct {
	PlanId         string
	Scope          string
	Page           int
	PerPage        int
	LogsPerPage    int
	Permanent      bool
}


func ParamsFromRequest(r *http.Request) *Params {
	params := &Params{}

	props := mux.Vars(r)
	query := r.URL.Query()


	if val, ok := props["plan_id"]; ok {
		params.PlanId = val
	}

	params.Scope = query.Get("scope")

	if val, err := strconv.Atoi(query.Get("page")); err != nil || val < 0 {
		params.Page = PAGE_DEFAULT
	} else {
		params.Page = val
	}

	if val, err := strconv.ParseBool(query.Get("permanent")); err == nil {
		params.Permanent = val
	}

	if val, err := strconv.Atoi(query.Get("per_page")); err != nil || val < 0 {
		params.PerPage = PER_PAGE_DEFAULT
	} else if val > PER_PAGE_MAXIMUM {
		params.PerPage = PER_PAGE_MAXIMUM
	} else {
		params.PerPage = val
	}

	if val, err := strconv.Atoi(query.Get("logs_per_page")); err != nil || val < 0 {
		params.LogsPerPage = LOGS_PER_PAGE_DEFAULT
	} else if val > LOGS_PER_PAGE_MAXIMUM {
		params.LogsPerPage = LOGS_PER_PAGE_MAXIMUM
	} else {
		params.LogsPerPage = val
	}

	return params
}
