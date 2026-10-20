package features_users_transport

import (
	"fmt"
	"net/http"

	core_logger "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/logger"
	"github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/transport/http/response"
)

func (h *UserHTTPHandler) LogoutUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	ResponseHandler := response.NewHandlerResponse(log, w)

	fmt.Println(123213)

	sessionID, err := r.Cookie("sessionID")
	fmt.Println(sessionID.Value)
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
