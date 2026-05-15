package feature_transport_session

import (
	"fmt"
	"net/http"

	core_logger "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/logger"
	"github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/transport/http/response"
	core_http_utils "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/transport/http/utils"
)

// GetSessionData   godoc
// @Summary      	Get cookie data
// @Description  	Get cookie by sessionID
// @Tags         	session
// @Produce      	json
// @Param        	sessionID  path string true         "ID for get something info"
// @Success 201		{object} core_domain.CookieData		"session cookieData return"
// @Failure      	500 {object} response.ErrorResponse "Internal server error"
// @Router       	/session/{sessionID}	[get]
func (h *SessionHandler) GetSessionData(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	ResponseHandler := response.NewHandlerResponse(log, w)

	sessionID, err := core_http_utils.GetValuePathString(r, "sessionID")
	if err != nil {
		ResponseHandler.ErrorResponse(err, fmt.Sprintf("Invalid session ID: %v", sessionID))
		return
	}

	data, err := h.SessionService.GetSessionData(ctx, sessionID)
	if err != nil {
		ResponseHandler.ErrorResponse(err, fmt.Sprintf("Error getting session data: %v", sessionID))
		return
	}

	ResponseHandler.JSONResponseHandler(http.StatusOK, data)
}
