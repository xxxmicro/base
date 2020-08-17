package config

import(
	"github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/config"
	"github.com/micro/go-micro/v2/config/source"
)

func NewConfigProvider(source source.Source) config.Config {
	cfg, _ := config.NewConfig()
	if err := cfg.Load(source); err != nil {
		logger.Log(logger.ErrorLevel, err)
	}

	return cfg
}