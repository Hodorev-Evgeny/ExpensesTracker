package features_users_transport

import (
	"net/http"

	core_logger "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/logger"
	"github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/transport/http/response"
)

func (h *UserHTTPHandler) LogoutUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	ResponseHandler := response.NewHandlerResponse(log, w)

	sessionID, err := r.Cookie("sessionID")
	if err != nil {
		ResponseHandler.ErrorResponse(err, "error getting sessionID")
		return
	}

	err = h.userService.LogoutUser(ctx, sessionID.Value)
	if err != nil {
		ResponseHandler.ErrorResponse(err, "error logouting user")
		return
	}

	ResponseHandler.JSONResponseHandler(http.StatusNoContent, nil)
}
