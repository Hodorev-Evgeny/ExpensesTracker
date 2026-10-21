package feature_transport_session

import (
	"fmt"
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
// @Param        	sessionID  path string true         "ID for get something info"
// @Success 201		{object} core_domain.CookieData		"session cookieData return"
// @Failure      	500 {object} response.ErrorResponse "Internal server error"
// @Router       	/users/me	[get]
func (h *SessionHandler) GetSessionData(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	ResponseHandler := response.NewHandlerResponse(log, w)

	userID, ok := ctx.Value("userID").(int)
	if !ok {
		ResponseHandler.ErrorResponse(core_errors.ErrorBadRequest, fmt.Sprintf("Error getting userId"))
		return
	}

	fmt.Println(userID, 111111111111)
	data := core_domain.CookieData{UserID: userID}

	ResponseHandler.JSONResponseHandler(http.StatusOK, data)
}
