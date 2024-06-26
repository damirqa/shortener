package app

import (
	"fmt"
	"github.com/damirqa/shortener/internal/infrastructure/logger"
	"github.com/gorilla/mux"
	"log"
	"net/http"

	"github.com/damirqa/shortener/cmd/config"
	URLDomainLocalRepository "github.com/damirqa/shortener/internal/domain/url/repository/local"
	URLDomainService "github.com/damirqa/shortener/internal/domain/url/service"
	URLUseCase "github.com/damirqa/shortener/internal/usecase/url"

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

func (app *App) initUseCases() {
	app.UseCases = &usecase.UseCases{
		URLUseCase: app.URLUseCase,
	}
}

func (app *App) initHTTPServer() {
	router := mux.NewRouter()
	router.Use(middleware.LogMW)
	handlers.RegisterHandlers(router, app.UseCases)

	app.httpServer = &http.Server{
		Addr:    config.Instance.GetAddress(),
		Handler: router,
	}
}

func (app *App) Start() {

	// http server
	{
		err := app.httpServer.ListenAndServe()
		if err != nil {
			fmt.Println(err.Error())
			panic(err)
		}
	}
}
