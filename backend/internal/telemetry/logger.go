package telemetry

import (
	"github.com/example/team-stack/backend/internal/config"
	"go.uber.org/zap"
)

func NewLogger(cfg *config.Config) (*zap.SugaredLogger, error) {
	var zcfg zap.Config
	if cfg.Env == "prod" {
		zcfg = zap.NewProductionConfig()
	} else {
		zcfg = zap.NewDevelopmentConfig()
	}

	if err := zcfg.Level.UnmarshalText([]byte(cfg.LogLevel)); err != nil {
		return nil, err
	}

	l, err := zcfg.Build()
	if err != nil {
		return nil, err
	}

	return l.Sugar(), nil
}
