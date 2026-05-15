package feature_service_web

import (
	"fmt"
	"os"
	"path"

	core_domain "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/domain"
)

func (s *WebService) GetLogin() (core_domain.File, error) {
	htmlPath := path.Join(
		os.Getenv("PROJECT_ROOT"),
		"public/login.html",
	)

	htmlFilePath, err := s.repository.GetFile(htmlPath)
	if err != nil {
		return core_domain.File{}, fmt.Errorf("error while getting register page file: %w", err)
	}

	return htmlFilePath, nil
}
