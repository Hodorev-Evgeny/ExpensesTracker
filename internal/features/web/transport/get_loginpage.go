package feature_transport_web

import (
	"net/http"

	core_logger "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/logger"
	"github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/transport/http/response"
)

func (h *WebTransport) GetLogin(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := response.NewHandlerResponse(log, w)

	log.Debug("start handling main page")
	htmlMainPage, err := h.service.GetLogin()
	if err != nil {
		responseHandler.ErrorResponse(err, "error getting main page")
		return
	}

	responseHandler.HTMLResponse(htmlMainPage.Buffer())
}
