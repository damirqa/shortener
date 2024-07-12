package app

import (
	"context"
	"database/sql"
	"errors"
	"github.com/damirqa/shortener/cmd/config"
	"github.com/damirqa/shortener/internal/domain/url/repository"
	URLDomainDBRepository "github.com/damirqa/shortener/internal/domain/url/repository/db"
	URLDomainLocalRepository "github.com/damirqa/shortener/internal/domain/url/repository/local"
	URLDomainService "github.com/damirqa/shortener/internal/domain/url/service"
	"github.com/damirqa/shortener/internal/infrastructure/logger"
	URLUseCase "github.com/damirqa/shortener/internal/usecase/url"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"log"
	"net/http"
	"time"

	"github.com/damirqa/shortener/internal/handlers"
	"github.com/damirqa/shortener/internal/middleware"
	"github.com/damirqa/shortener/internal/usecase"
)

type App struct {
	httpServer *http.Server

	// url
	URLDomainService *URLDomainService.URLService
	//URLDomainRepository *URLDomainLocalRepository.URLLocalRepository
	URLDomainRepository repository.URLRepository
	URLUseCase          *URLUseCase.UseCase

	// use cases
	UseCases *usecase.UseCases
}

func (app *App) Init() {
	app.initConfig()
	app.initLogger()
	app.initURL()
	app.recoveryURL()
	app.initDBSchemas()
	app.initUseCases()
	app.initHTTPServer()
}

func (app *App) initConfig() {
	config.Instance = config.Init()
}

func (app *App) initLogger() {
	if err := logger.Initialize(config.Instance.LogLevel); err != nil {
		log.Fatal(err)
	}
}

func (app *App) initURL() {
	URLDBRepository, err := URLDomainDBRepository.New()
	if err != nil {
		logger.GetLogger().Error("problem with connection to db", zap.Error(err))
		app.URLDomainRepository = URLDomainLocalRepository.New()
	} else {
		app.URLDomainRepository = URLDBRepository
	}

	app.URLDomainService = URLDomainService.New(app.URLDomainRepository)
	app.URLUseCase = URLUseCase.New(app.URLDomainService)
}

func (app *App) recoveryURL() {
	app.URLDomainService.LoadFromFile()
}

func (app *App) initDBSchemas() {
	db, err := sql.Open("pgx", config.Instance.DatabaseDSN)
	if err != nil {
		logger.GetLogger().Error("problem with connection to db", zap.Error(err))
		return
	}

	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			logger.GetLogger().Error("problem with closing connection", zap.Error(err))
		}
	}(db)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		logger.GetLogger().Error("db ping failed", zap.Error(err))
		return
	}

	query := `
	CREATE TABLE IF NOT EXISTS urls (
		short VARCHAR(255) PRIMARY KEY,
		long TEXT NOT NULL
	);`

	_, err = db.ExecContext(ctx, query)
	if err != nil {
		logger.GetLogger().Error("failed to create table", zap.Error(err))
		return
	}

	logger.GetLogger().Info("table created successfully")
}

func (app *App) initUseCases() {
	app.UseCases = &usecase.UseCases{
		URLUseCase: app.URLUseCase,
	}
}

func (app *App) initHTTPServer() {
	router := mux.NewRouter()
	router.Use(middleware.LogMiddleware)
	router.Use(middleware.GzipMiddleware)
	handlers.RegisterHandlers(router, app.UseCases)

	app.httpServer = &http.Server{
		Addr:    config.Instance.GetAddress(),
		Handler: router,
	}
}

func (app *App) Start() {
	if err := app.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.GetLogger().Fatal("Could not listen",
			zap.String("address", config.Instance.GetAddress()),
			zap.Error(err))
	}
}

func (app *App) Shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := app.httpServer.Shutdown(ctx); err != nil {
		logger.GetLogger().Fatal("Server forced to shutdown", zap.Error(err))
	}

	app.URLDomainService.SaveToFile()
	app.URLDomainRepository.Close()
}
