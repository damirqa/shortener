package handlers

import (
	"database/sql"
	"github.com/damirqa/shortener/cmd/config"
	"github.com/damirqa/shortener/internal/infrastructure/logger"
	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
	"net/http"
)

func Ping() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		db, err := sql.Open("pgx", config.Instance.DatabaseDSN)

		if err != nil {
			logger.GetLogger().Error("problem with connection to db", zap.Error(err))
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		err = db.Ping()
		if err != nil {
			logger.GetLogger().Error("db ping failed", zap.Error(err))
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		writer.WriteHeader(http.StatusOK)

		defer func(db *sql.DB) {
			err := db.Close()
			if err != nil {
				logger.GetLogger().Error("problem with db", zap.Error(err))
			}
		}(db)
	}
}
