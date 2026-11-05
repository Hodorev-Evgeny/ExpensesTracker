package features_users_transport

import (
	"fmt"
	"net/http"

	core_domain "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/domain"
	core_errors "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/errors"
	core_logger "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/logger"
	"github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/transport/http/response"
)

// ExtraditionUser       godoc
// @Summary      Получение пользователя
// @Description  Получение конкретного пользователя по его ID
// @Tags         users
// @Produce      json
// @Param        id  path int true                          "ID получаемого пользователя"
// @Success      200 {object} UserDTOResponse                  "Пользователь успешно найден"
// @Failure      400 {object} response.ErrorResponse "Bad request"
// @Failure      404 {object} response.ErrorResponse "User not found"
// @Failure      500 {object} response.ErrorResponse "Internal server error"
// @Router       /users/{id} [get]
func (h *UserHTTPHandler) ExtraditionUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := response.NewHandlerResponse(log, w)

	log.Info("End processing get value path")
	id, ok := core_domain.UserIDFromContext(ctx)
	if !ok {
		responseHandler.ErrorResponse(
			fmt.Errorf("user not found in context: %w", core_errors.ErrorUnauthorized),
			"user is not registered")
		return
	}

	log.Info("Start processing get user")
	domainUser, err := h.userService.ExtraditionUser(ctx, id)
	if err != nil {
		responseHandler.ErrorResponse(err, "error try extradition user")
		return
	}
	log.Info("End processing get user")

	req := DomainFromResponse(domainUser)

	responseHandler.JSONResponseHandler(http.StatusOK, req)
}
