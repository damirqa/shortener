package api

import (
	"encoding/json"
	"github.com/damirqa/shortener/internal/domain/url/model"
	"github.com/damirqa/shortener/internal/infrastructure/logger"
	"github.com/damirqa/shortener/internal/middleware"
	URLUseCase "github.com/damirqa/shortener/internal/usecase/url"
	"go.uber.org/zap"
	"io"
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
		writer.WriteHeader(http.StatusOK)

		_, _ = writer.Write(resp)
	}
}

type DeleteURLsRequest []string

func DeleteUserLinks(useCase URLUseCase.UseCaseInterface) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		userID, ok := request.Context().Value(middleware.UserIDKey).(string)
		if !ok {
			http.Error(writer, "Not authorized", http.StatusUnauthorized)
			return
		}

		var deleteURLsRequest DeleteURLsRequest
		if err := json.NewDecoder(request.Body).Decode(&deleteURLsRequest); err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}

		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				logger.GetLogger().Error(err.Error())
			}
		}(request.Body)

		delChan := make(chan error, 1)
		go func() {
			delChan <- useCase.DeleteUserLinks(userID, deleteURLsRequest)
		}()

		writer.WriteHeader(http.StatusAccepted)

		go func() {
			if err := <-delChan; err != nil {
				logger.GetLogger().Error("Error deleting user links", zap.Error(err))
			}
		}()
	}
}
