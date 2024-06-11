package app

import (
	"fmt"
	"github.com/damirqa/shortener/cmd/config"
	"github.com/damirqa/shortener/internal/handlers"
	"github.com/damirqa/shortener/internal/usecase"
	"github.com/gorilla/mux"
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

		for {
			if isPortAvailable(address) {
				break
			}
			time.Sleep(2 * time.Second) // Ждать 2 секунды перед следующей проверкой
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

func isPortAvailable(address string) bool {
	conn, err := net.DialTimeout("tcp", address, 1*time.Second)
	if err != nil {
		return true // Порт доступен
	}
	conn.Close()
	return false // Порт занят
}
