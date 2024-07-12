package api

import (
	"encoding/json"
	URLModels "github.com/damirqa/shortener/internal/domain/url/model"
	"github.com/damirqa/shortener/internal/infrastructure/logger"
	URLUseCase "github.com/damirqa/shortener/internal/usecase/url"
	"io"
	"net/http"
)

type URLRequest struct {
	Link string `json:"url"`
}

type URLResponse struct {
	Link string `json:"result"`
}

func ShortenURL(useCase URLUseCase.UseCaseInterface) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		var urlRequest URLRequest
		if err := json.NewDecoder(request.Body).Decode(&urlRequest); err != nil {
			// todo: как залогировать request.Body, ведь это же поток
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}

		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				logger.GetLogger().Error(err.Error())
			}
		}(request.Body)

		shortURL := useCase.Generate(urlRequest.Link)

		urlResponse := URLResponse{Link: string(shortURL)}

		resp, err := json.Marshal(urlResponse)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
		}

		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusCreated)

		_, err = writer.Write(resp)
		if err != nil {
			http.Error(writer, "Error generate short url", http.StatusInternalServerError)
			return
		}
	}
}

func ShortenURLSBatch(useCase URLUseCase.UseCaseInterface) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		var urlsRequest []URLModels.URLRequestWithCorrelationId
		if err := json.NewDecoder(request.Body).Decode(&urlsRequest); err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}

		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				logger.GetLogger().Error(err.Error())
			}
		}(request.Body)

		shortURLs, err := useCase.GenerateBatch(urlsRequest)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		resp, err := json.Marshal(shortURLs)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusCreated)

		_, err = writer.Write(resp)
		if err != nil {
			http.Error(writer, "Error generate short urls", http.StatusInternalServerError)
			return
		}
	}
}
