package features_users_transport

import (
	"fmt"
	"net/http"

	core_domain "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/domain"
	core_errors "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/errors"
	core_logger "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/logger"
	"github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/transport/http/response"
	core_http_query_parm "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/transport/http/utils"
)

type GetUsersResponse []UserDTOResponse

// GetUsers      godoc
// @Summary      Список пользователей
// @Description  Просмотр списка пользователей с опциональной пагинацией
// @Tags         users
// @Produce      json
// @Param        limit  query int false                        "Размер страницы с пользователями"
// @Param        offset query int false                        "Смещение страницы с пользователями"
// @Success      200 {object} GetUsersResponse                 "Успешное получение списка пользователей"
// @Failure      400 {object} response.ErrorResponse "Bad request"
// @Failure      500 {object} response.ErrorResponse "Internal server error"
// @Router       /users [get]
func (h *UserHTTPHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	ResponseHandler := response.NewHandlerResponse(log, w)

	log.Info("Start processing get query parameters")
	limit, offset, err := GetLimitAnsOffset(r)
	if err != nil {
		ResponseHandler.ErrorResponse(err, "error parsing limit/offset")
		return
	}

	log.Info("Start processing get users")
	usersDomain, err := h.userService.GetUsers(ctx, limit, offset)
	if err != nil {
		ResponseHandler.ErrorResponse(err, "error getting users")
		return
	}
	log.Info("End processing get users")

	userRsponse := DomainFromDTOesponse(usersDomain)

	ResponseHandler.JSONResponseHandler(http.StatusOK, userRsponse)
}

func GetLimitAnsOffset(r *http.Request) (*int, *int, error) {
	limit, err := core_http_query_parm.GetIntQueryParm(r, "limit")

	if err != nil {
		return nil, nil, fmt.Errorf("error parsing limit: %w", core_errors.ErrorValidation)
	}

	offset, err := core_http_query_parm.GetIntQueryParm(r, "offset")
	if err != nil {
		return nil, nil, fmt.Errorf("error parsing offset: %w", core_errors.ErrorValidation)
	}

	return limit, offset, nil
}

func DomainFromDTOesponse(rows []core_domain.User) GetUsersResponse {
	users := make([]UserDTOResponse, 0)
	for _, user := range rows {
		userResponse := DomainFromResponse(user)
		users = append(users, userResponse)
	}
	return users
}
