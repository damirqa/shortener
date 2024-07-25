package api

import (
	"encoding/json"
	"github.com/damirqa/shortener/internal/domain/url/model"
	"github.com/damirqa/shortener/internal/infrastructure/logger"
	"github.com/damirqa/shortener/internal/middleware"
	URLUseCase "github.com/damirqa/shortener/internal/usecase/url"
	"net/http"
)

func GetAllUserLinks(useCase URLUseCase.UseCaseInterface) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		userID, ok := request.Context().Value(middleware.UserIDKey).(string)
		if !ok {
			http.Error(writer, "Forbidden", http.StatusForbidden)
			return
		}

		urls, err := useCase.GetAllUserLinks(userID)
		if err != nil {
			http.Error(writer, "Internal server error", http.StatusInternalServerError)
			logger.GetLogger().Error(err.Error())
			return
		}

		userURLsResponse := make([]model.URLResponse, 0, len(urls))
		for _, u := range urls {
			userURLsResponse = append(userURLsResponse, model.URLResponse{ShortURL: u.ShortURL, OriginalURL: u.OriginalURL})
		}

		resp, err := json.Marshal(userURLsResponse)
		if err != nil {
			http.Error(writer, "Internal server error", http.StatusInternalServerError)
			logger.GetLogger().Error(err.Error())
			return
		}

		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusCreated)

		_, _ = writer.Write(resp)
	}
}
