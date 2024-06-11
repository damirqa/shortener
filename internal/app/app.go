package app

import (
	"fmt"
	"github.com/damirqa/shortener/cmd/config"
	"github.com/damirqa/shortener/internal/handlers"
	"github.com/damirqa/shortener/internal/usecase"
	"github.com/gorilla/mux"
	"log"
	"net"
	"net/http"
	"time"

	URLDomainLocalRepository "github.com/damirqa/shortener/internal/domain/url/repository/local"
	URLDomainService "github.com/damirqa/shortener/internal/domain/url/service"
	URLUseCase "github.com/damirqa/shortener/internal/usecase/url"
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
		config.ParseFlags()
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
		address := fmt.Sprintf("%s:%s", config.FlagRunAddr, config.FlagRunPort)

		timeout := 10 * time.Second
		conn, err := net.DialTimeout("tcp", address, timeout)
		if err != nil {
			log.Printf("Port %s is available, starting server...\n", config.FlagRunPort)
		} else {
			err := conn.Close()
			if err != nil {
				_ = fmt.Errorf(err.Error())
			}
			log.Fatalf("Port %s is already in use.\n", config.FlagRunPort)
		}

		router := mux.NewRouter()
		handlers.RegisterHandlers(router, app.UseCases)

		app.httpServer = &http.Server{
			Addr:    config.FlagRunAddr + config.FlagRunPort,
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
