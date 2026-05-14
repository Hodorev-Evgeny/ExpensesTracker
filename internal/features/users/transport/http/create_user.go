package features_users_transport

import (
	"encoding/json"
	"net/http"

	core_logger "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/logger"
	"github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/transport/http/response"
)

type CreateUserRequest struct {
	FullName string  `json:"full_name" validate:"required,min=3,max=32" example:"John Doe"`
	Password string  `json:"password" validate:"required,min=8,max=32" example:"123456"`
	Email    string  `json:"email" validate:"required,email" example:"example@gmail.com"`
	Phone    *string `json:"phone" validate:"required,min=11,max=11, startswith=+" example:"+71124234312"`
}

type CreateUserResponse UserDTOResponse

// CreateUser   godoc
// @Summary     Создать пользователя
// @Description Создать нового пользователя в системе
// @Tags        users
// @Accept      json
// @Produce     json
// @Param       request body     CreateUserRequest  true "CreateUser тело запроса"
// @Success     201     {object} CreateUserResponse "Успешно созданный пользователь"
// @Failure     400     {object} response.ErrorResponse "Bad request"
// @Failure     500     {object} response.ErrorResponse "Internal server error"
// @Router      /users [post]
func (h *UserHTTPHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	RsponceHandler := response.NewHandlerResponse(log, w)

	log.Info("Start decoding body")
	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RsponceHandler.ErrorResponse(err, "error decoding request body")
		return
	}

	log.Info("Start processing creating new user")
	userDomain, err := h.userService.CreateUser(ctx, DTOFromDomain(req))
	if err != nil {
		RsponceHandler.ErrorResponse(err, "error creating user")
		return
	}
	log.Info("End processing creating new user")

	sessionID, er := h.userService.CreateCache(ctx, userDomain)
	if er != nil {
		RsponceHandler.ErrorResponse(er, "error creating string session id")
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:  "sessionID",
		Value: sessionID,
		Path:  "/",
	})

	log.Info("Start writing response")
	respons := DomainFromResponse(userDomain)
	if err := json.NewEncoder(w).Encode(respons); err != nil {
		RsponceHandler.ErrorResponse(err, "error writing response")
		return
	}
	w.WriteHeader(http.StatusCreated)
}
