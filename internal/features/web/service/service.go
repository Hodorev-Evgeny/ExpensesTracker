package feature_service_web

import core_domain "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/domain"

type WebRepository interface {
	GetFile(
		htmlPath string,
	) (core_domain.File, error)
}

type WebService struct {
	repository WebRepository
}

func NewWebService(
	repository WebRepository,
) *WebService {
	return &WebService{
		repository: repository,
	}
}
