package feature_repository_file_system

import (
	"errors"
	"fmt"
	"os"

	core_domain "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/domain"
	core_errors "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/errors"
)

func (r *WebRepository) GetFile(filepath string) (core_domain.File, error) {
	bytes, err := os.ReadFile(filepath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return core_domain.File{}, fmt.Errorf("file does not exist: %w", core_errors.ErrorNotFoud)
		}

		return core_domain.File{}, fmt.Errorf("file does not exist: %w", core_errors.ErrorNotFoud)
	}

	file := core_domain.NewFile(bytes)

	return file, nil
}
