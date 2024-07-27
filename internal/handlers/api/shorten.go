package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/damirqa/shortener/cmd/config"
	URLModels "github.com/damirqa/shortener/internal/domain/url/model"
	dberror "github.com/damirqa/shortener/internal/error"
	"github.com/damirqa/shortener/internal/infrastructure/logger"
	"github.com/damirqa/shortener/internal/middleware"
	URLUseCase "github.com/damirqa/shortener/internal/usecase/url"
	"go.uber.org/zap"
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

		var bodyBuffer bytes.Buffer
		reader := io.TeeReader(request.Body, &bodyBuffer)

		if err := json.NewDecoder(reader).Decode(&urlRequest); err != nil {
			logger.GetLogger().Error("Error decoding request body: %v. Body: %s", zap.Error(err), zap.String("body", bodyBuffer.String()))
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}

		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				logger.GetLogger().Error(err.Error())
			}
		}(request.Body)

		userID := request.Context().Value(middleware.UserIDKey).(string)

		URLEntity, err := useCase.Generate(urlRequest.Link, userID)
		if err != nil {
			var uniqueErr *dberror.UniqueConstraintError
			if errors.As(err, &uniqueErr) {
				writer.Header().Set("Content-Type", "application/json")
				writer.WriteHeader(http.StatusConflict)
			} else {
				logger.GetLogger().Error(err.Error())
				http.Error(writer, "Error generate short url", http.StatusInternalServerError)
				return
			}
		}

		fullURL := config.Instance.GetResultAddress() + "/" + URLEntity.ShortURL
		urlResponse := URLResponse{Link: fullURL}

		resp, err := json.Marshal(urlResponse)
		if err != nil {
			logger.GetLogger().Error(err.Error())
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusCreated)

		_, _ = writer.Write(resp)
	}
}

func ShortenURLSBatch(useCase URLUseCase.UseCaseInterface) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		var urlsRequest []URLModels.URLRequestWithCorrelationID
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

		URLEntities, err := useCase.GenerateBatch(urlsRequest)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		var URLSResponseWithCorrelationID []URLModels.URLResponseWithCorrelationID
		for _, url := range URLEntities {
			URLResponseWithCorrelationID := URLModels.URLResponseWithCorrelationID{CorrelationID: url.CorrelationID, ShortURL: config.Instance.GetResultAddress() + "/" + url.ShortURL}
			URLSResponseWithCorrelationID = append(URLSResponseWithCorrelationID, URLResponseWithCorrelationID)
		}

		resp, err := json.Marshal(URLSResponseWithCorrelationID)
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
