package app

import (
	"fmt"
	"github.com/damirqa/shortener/cmd/config"
	URLDomainLocalRepository "github.com/damirqa/shortener/internal/domain/url/repository/local"
	URLDomainService "github.com/damirqa/shortener/internal/domain/url/service"
	"github.com/damirqa/shortener/internal/handlers"
	"github.com/damirqa/shortener/internal/usecase"
	URLUseCase "github.com/damirqa/shortener/internal/usecase/url"
	"github.com/gorilla/mux"
	"net/http"
)

type App struct {
	httpServer *http.Server

	// url
	URLDomainService    *URLDomainService.Service
	URLDomainRepository *URLDomainLocalRepository.Local
	URLUseCase          *URLUseCase.UseCase

	// use cases
	UseCases *usecase.UseCases
}

func (app *App) Init() {
	//config
	{
		config.ConfigInstance = config.Init()
	}

	// url
	{
		app.URLDomainRepository = URLDomainLocalRepository.New()
		app.URLDomainService = URLDomainService.New(app.URLDomainRepository)
		app.URLUseCase = URLUseCase.New(app.URLDomainService)
	}

	// use cases
	{
		app.UseCases = &usecase.UseCases{
			URLUseCase: app.URLUseCase,
		}
	}

	// http server
	{
		router := mux.NewRouter()
		handlers.RegisterHandlers(router, app.UseCases)

		app.httpServer = &http.Server{
			Addr:    config.ConfigInstance.GetFullAddress(),
			Handler: router,
		}
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
