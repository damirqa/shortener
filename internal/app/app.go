package app

import (
	"context"
	"errors"
	"github.com/damirqa/shortener/cmd/config"
	URLDomainLocalRepository "github.com/damirqa/shortener/internal/domain/url/repository/local"
	URLDomainService "github.com/damirqa/shortener/internal/domain/url/service"
	"github.com/damirqa/shortener/internal/infrastructure/logger"
	URLUseCase "github.com/damirqa/shortener/internal/usecase/url"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/damirqa/shortener/internal/handlers"
	"github.com/damirqa/shortener/internal/middleware"
	"github.com/damirqa/shortener/internal/usecase"
)

type App struct {
	httpServer *http.Server

	// url
	URLDomainService    *URLDomainService.URLService
	URLDomainRepository *URLDomainLocalRepository.URLLocalRepository
	URLUseCase          *URLUseCase.UseCase

	// use cases
	UseCases *usecase.UseCases
}

func (app *App) Init() {
	app.initConfig()
	app.initLogger()
	app.initURL()
	app.recoveryURL()
	app.initUseCases()
	app.initHTTPServer()
}

func (app *App) initConfig() {
	config.Instance = config.Init()
}

func (app *App) initLogger() {
	if err := logger.Initialize(config.Instance.FlagLogLevel); err != nil {
		log.Fatal(err)
	}
}

func (app *App) initURL() {
	app.URLDomainRepository = URLDomainLocalRepository.New()
	app.URLDomainService = URLDomainService.New(app.URLDomainRepository)
	app.URLUseCase = URLUseCase.New(app.URLDomainService)
}

func (app *App) recoveryURL() {
	app.URLDomainService.RecoveryURLsFromFile()
}

func (app *App) initUseCases() {
	app.UseCases = &usecase.UseCases{
		URLUseCase: app.URLUseCase,
	}
}

func (app *App) initHTTPServer() {
	router := mux.NewRouter()
	router.Use(middleware.LogMW)
	router.Use(middleware.GzipMW)
	handlers.RegisterHandlers(router, app.UseCases)

	app.httpServer = &http.Server{
		Addr:    config.Instance.GetAddress(),
		Handler: router,
	}
}

func (app *App) Start() {
	if err := app.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("Could not listen on %s: %v\n", config.Instance.GetAddress(), err)
	}
}

func (app *App) Listen() {
	signalCtx, signalCtxCancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer signalCtxCancel()

	// wait signal
	<-signalCtx.Done()
}

func (app *App) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := app.httpServer.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	app.URLDomainService.DumpToFile()
}
