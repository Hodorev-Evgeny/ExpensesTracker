package features_users_transport

import (
	"net/http"

	core_domain "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/domain"
	core_logger "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/logger"
	"github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/transport/http/response"
	core_http_utils "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/transport/http/utils"
)

type LoginHandler struct {
	FullName string `json:"full_name"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type StatusRequest struct {
	Status string `json:"status"`
}

// LoginUser     godoc
// @Summary      Create session
// @Description  Create session when user register
// @Tags         users
// @Produce      json
// @Param        request body LoginHandler true      "This Body need to login user"
// @Success 201		 {object} StatusRequest		 	 "session created"
// @Failure      400 {object} response.ErrorResponse "Bad request"
// @Failure      500 {object} response.ErrorResponse "Internal server error"
// @Router       /users/login [post]
func (h *UserHTTPHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	ResponseHandler := response.NewHandlerResponse(log, w)

	var user LoginHandler
	if err := core_http_utils.DecodeJSON(&user, r); err != nil {
		ResponseHandler.ErrorResponse(err, "error decoding body")
		return
	}

	loginUser := core_domain.CreateLoginUser(user.FullName, user.Email, user.Password)
	sessionID, err := h.userService.LoginUser(ctx, loginUser)
	if err != nil {
		ResponseHandler.ErrorResponse(err, "error logging in")
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:  "sessionID",
		Value: sessionID,
	})

	req := StatusRequest{
		Status: "session created",
	}

	ResponseHandler.JSONResponseHandler(http.StatusCreated, req)
}
