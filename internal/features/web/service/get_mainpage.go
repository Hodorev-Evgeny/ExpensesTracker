package feature_service_web

import (
	"fmt"
	"os"
	"path"

	core_domain "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/domain"
)

func (s *WebService) GetMainPage() (core_domain.File, error) {
	htmlPath := path.Join(
		os.Getenv("PROJECT_ROOT"),
		"public/index.html",
	)

	htmlFilePath, err := s.repository.GetFile(htmlPath)
	if err != nil {
		return core_domain.File{}, fmt.Errorf("error while getting main page file: %w", err)
	}

	return htmlFilePath, nil
}
