package app

import (
	"fmt"
	"github.com/damirqa/shortener/internal/handlers"
	"github.com/damirqa/shortener/internal/usecase"
	"github.com/gorilla/mux"
	"net/http"

	urlDomainLocalRepository "github.com/damirqa/shortener/internal/domain/url/repository/local"
	urlDomainService "github.com/damirqa/shortener/internal/domain/url/service"
	urlUseCase "github.com/damirqa/shortener/internal/usecase/url"
)

type App struct {
	httpServer *http.Server

	// url
	UrlDomainService    *urlDomainService.Service
	UrlDomainRepository *urlDomainLocalRepository.Local
	UrlUseCase          *urlUseCase.UseCase

	// use cases
	UseCases *usecase.UseCases
}

func (app *App) Init() {
	// url
	{
		app.UrlDomainRepository = urlDomainLocalRepository.New()
		app.UrlDomainService = urlDomainService.New(app.UrlDomainRepository)
		app.UrlUseCase = urlUseCase.New(app.UrlDomainService)
	}

	// use cases
	{
		app.UseCases = &usecase.UseCases{
			URLUseCase: app.UrlUseCase,
		}
	}

	// http server
	{
		router := mux.NewRouter()
		handlers.RegisterHandlers(router, app.UseCases)

		app.httpServer = &http.Server{
			Addr:    "127.0.0.1:8080",
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
