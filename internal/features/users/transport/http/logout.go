package features_users_transport

import (
	"net/http"

	core_errors "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/errors"
	core_logger "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/logger"
	"github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/transport/http/response"
	core_http_utils "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/transport/http/utils"
)

func (h *UserHTTPHandler) LogoutUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	ResponseHandler := response.NewHandlerResponse(log, w)

	sessionID, err := r.Cookie("sessionID")
	if err != nil {
		ResponseHandler.ErrorResponse(core_errors.ErrorUnauthorized, "error getting sessionID")
		return
	}

	err = h.userService.LogoutUser(ctx, sessionID.Value)
	if err != nil {
		ResponseHandler.ErrorResponse(err, "error logouting user")
		return
	}

	core_http_utils.DeleteSessionCookie(w)

	ResponseHandler.JSONResponseHandler(http.StatusNoContent, nil)
}
