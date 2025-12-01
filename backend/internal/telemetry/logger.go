package telemetry

import (
    "github.com/example/team-stack/backend/internal/config"
    "go.uber.org/zap"
)

func NewLogger(cfg *config.Config) *zap.SugaredLogger {
    var zcfg zap.Config
    if cfg.Env == "prod" {
        zcfg = zap.NewProductionConfig()
    } else {
        zcfg = zap.NewDevelopmentConfig()
    }
    zcfg.Level.UnmarshalText([]byte(cfg.LogLevel))
    l, _ := zcfg.Build()
    return l.Sugar()
}
