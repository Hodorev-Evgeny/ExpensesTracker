package features_users_transport

import (
	"fmt"
	"net/http"

	core_domain "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/domain"
	core_errors "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/errors"
	core_logger "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/logger"
	"github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/transport/http/response"
)

// DeleteUser    godoc
// @Summary      Удаление пользователя
// @Description  Удаление существующего в системе пользователя по его ID
// @Tags         users
// @Param        id  path int true                          "ID удаляемого пользователя"
// @Success      204                                           "Успешное удаление пользователя"
// @Failure      400 {object} response.ErrorResponse "Bad request"
// @Failure      404 {object} response.ErrorResponse "User not found"
// @Failure      500 {object} response.ErrorResponse "Internal server error"
// @Router       /users/{id} [delete]
func (h *UserHTTPHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := response.NewHandlerResponse(log, w)

	log.Info("Start processing get value path")
	id, ok := core_domain.UserIDFromContext(ctx)
	if !ok {
		responseHandler.ErrorResponse(
			fmt.Errorf("user not found in context: %w", core_errors.ErrorUnauthorized),
			"user is not registered")
		return
	}

	log.Info("End start processing delete user")
	if err := h.userService.DeleteUser(ctx, id); err != nil {
		responseHandler.ErrorResponse(err, err.Error())
		return
	}
	log.Info("End processing delete user")

	responseHandler.JSONResponseHandler(http.StatusNoContent, nil)
}
