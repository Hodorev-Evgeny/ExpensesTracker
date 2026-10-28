package feature_transport_session

import (
	"net/http"

	core_domain "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/domain"
	core_errors "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/errors"
	core_logger "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/logger"
	"github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/transport/http/response"
)

// GetSessionData   godoc
// @Summary      	Get cookie data
// @Description  	Get cookie by sessionID
// @Tags         	session
// @Produce      	json
// @Success 200		{object} core_domain.CookieData		"session cookieData return"
// @Failure      	500 {object} response.ErrorResponse "Internal server error"
// @Router       	/users/me	[get]
func (h *SessionHandler) GetSessionData(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	ResponseHandler := response.NewHandlerResponse(log, w)

	userID, ok := core_domain.UserIDFromContext(ctx)
	if !ok {
		ResponseHandler.ErrorResponse(core_errors.ErrorUnauthorized, "userID not found in context")
		return
	}

	data := core_domain.CookieData{UserID: userID}

	ResponseHandler.JSONResponseHandler(http.StatusOK, data)
}
