package app

import (
	"github.com/OhBonsai/yogo/model"
	"github.com/OhBonsai/yogo/utils"
)

func (a *App) Config() *model.Config {
	if cfg := a.config.Load(); cfg != nil {
		return cfg.(*model.Config)
	}
	return &model.Config{}
}


func (a *App) LoadConfig(configFile string) *model.AppError {
	cfg, configPath, envConfig, err := utils.LoadConfig(configFile)
	if err != nil {
		return err
	}

	a.configFile = configPath
	a.config.Store(cfg)
	a.envConfig = envConfig

	return nil
}