package usecase

import (
	URLUseCase "github.com/damirqa/shortener/internal/usecase/url"
)

type UseCases struct {
	URLUseCase URLUseCase.UseCaseInterface
}
