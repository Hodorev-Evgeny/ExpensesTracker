package features_users_transport

import (
	"errors"
	"net/http"

	core_domain "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/domain"
	core_logger "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/logger"
	"github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/transport/http/response"
	core_http_types "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/transport/http/types"
	core_http_utils "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/transport/http/utils"
	"go.uber.org/zap"
)

type RequestPatchUser struct {
	FullName core_http_types.Nullable[string] `json:"full_name" swaggertype:"string" example:"John Doe"`
	Email    core_http_types.Nullable[string] `json:"email" swaggertype:"string" example:"example@gmai;.com"`
	Phone    core_http_types.Nullable[string] `json:"phone" swaggertype:"string" example:"555-555-5555"`
}

func (u *RequestPatchUser) Validate() error {
	if u.FullName.Set {
		if u.FullName.Value == nil {
			return errors.New("full name must not be empty")
		}
	}
	if u.Email.Set {
		if u.Email.Value == nil {
			return errors.New("email must not be empty")
		}
	}

	return nil
}

func CreateUserPatch(user RequestPatchUser) core_domain.UserPatch {
	return core_domain.NewUserPatch(user.FullName.ToDomain(),
		user.Email.ToDomain(),
		user.Phone.ToDomain(),
	)
}

// PatchUser     godoc
// @Summary      Изменение пользователя
// @Description Changing information about a user already existing in the system
// @Description ### Logic of updating fields (Three-state logic):
// @Description 1. **The field was not passed**: `phone_number` is ignored, the value in the database does not change
// @Description 2. **The value** is explicitly passed: `"phone_number": "+711122233344"` - sets a new phone number in the database
// @Description 3. **null passed**: `"phone_number": null` - clears the field in the database (set to NULL)
// @Description Restrictions: `full_name` cannot be set as null
// @Tags         users
// @Accept       json
// @Produce      json
// @Param		id path int true "ID of the user being modified"
// @Param 		request body RequestPatchUser true "PatchUser request body"
// @Success 	200 {object} UserDTOResponse "Successfully changed user"
// @Failure      400 {object} response.ErrorResponse "Bad request"
// @Failure      404 {object} response.ErrorResponse "User not found"
// @Failure      409 {object} response.ErrorResponse "Conflict"
// @Failure      500 {object} response.ErrorResponse "Internal server error"
// @Router       /users/{id} [patch]
func (h *UserHTTPHandler) PatchUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	ResponseHandler := response.NewHandlerResponse(log, w)

	log.Info("Start processing decoding body")
	var data RequestPatchUser
	if err := core_http_utils.DecodeJSON(&data, r); err != nil {
		ResponseHandler.ErrorResponse(err, "error parsing request")

		return
	}

	log.Info("Start processing patch user")
	userId, err := core_http_utils.GetValuePathInt(r, "id")
	if err != nil {
		ResponseHandler.ErrorResponse(err, "error parsing id")
		return
	}

	log.Debug("PatchUser",
		zap.Any("data", data),
	)

	userPatch := CreateUserPatch(data)
	log.Info("End processing patch user")

	log.Info("Start processing response body")
	userDomain, err := h.userService.PatchUser(ctx, userId, userPatch)
	if err != nil {
		ResponseHandler.ErrorResponse(err, "error patching user")
		return
	}

	ResponseHandler.JSONResponseHandler(http.StatusOK, userDomain)
}
