package web

import (
	"github.com/OhBonsai/yogo/app"
	"github.com/OhBonsai/yogo/mlog"
	"github.com/OhBonsai/yogo/model"
	"github.com/OhBonsai/yogo/utils"

	"net/http"

	goi18n "github.com/nicksnyder/go-i18n/i18n"
)

type Context struct {
	App *app.App
	Log *mlog.Logger
	Session model.Session
	Params        *Params
	Err           *model.AppError
	T             goi18n.TranslateFunc
	RequestId     string
	IpAddress     string
	Path          string
	siteURLHeader string
}


func (c *Context) LogError(err *model.AppError) {
	c.Log.Error(
		err.SystemMessage(utils.TDefault),
		mlog.String("err_where", err.Where),
		mlog.Int("http_code", err.StatusCode),
		mlog.String("err_details", err.DetailedError),
	)
}

func (c *Context) SetInvalidParam(parameter string) {
	c.Err = NewInvalidParamError(parameter)
}


func (c *Context) LogInfo(err *model.AppError) {
	// Filter out 401s
	if err.StatusCode == http.StatusUnauthorized {
		c.LogDebug(err)
	} else {
		c.Log.Info(
			err.SystemMessage(utils.TDefault),
			mlog.String("err_where", err.Where),
			mlog.Int("http_code", err.StatusCode),
			mlog.String("err_details", err.DetailedError),
		)
	}
}

func (c *Context) LogDebug(err *model.AppError) {
	c.Log.Debug(
		err.SystemMessage(utils.TDefault),
		mlog.String("err_where", err.Where),
		mlog.Int("http_code", err.StatusCode),
		mlog.String("err_details", err.DetailedError),
	)
}

func NewInvalidParamError(parameter string) *model.AppError {
	err := model.NewAppError("Context", "api.context.invalid_body_param.app_error", map[string]interface{}{"Name": parameter}, "", http.StatusBadRequest)
	return err
}
