package usecase

import (
	urlUseCase "github.com/damirqa/shortener/internal/usecase/url"
)

type UseCases struct {
	URLUseCase urlUseCase.ServiceInterface
}
