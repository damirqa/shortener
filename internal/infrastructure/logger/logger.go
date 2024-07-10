package logger

import (
	"go.uber.org/zap"
)

var log *zap.Logger = zap.NewNop()

func Initialize(level string) error {
	lvl, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return err
	}

	cfg := zap.NewProductionConfig()

	cfg.Level = lvl

	zl, err := cfg.Build()
	if err != nil {
		return err
	}

	log = zl

	return nil
}

func GetLogger() *zap.Logger {
	return log
}
