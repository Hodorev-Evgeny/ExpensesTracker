package feature_transport_static

import (
	"fmt"
	"net/http"

	core_domain "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/domain"
	core_logger "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/logger"
	"github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/transport/http/response"
	core_http_utils "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/transport/http/utils"
)

type StaticResponse struct {
	SumIncome      int                `json:"sum_income" example:"1255"`
	SumExpenditure int                `json:"sum_expenditure" example:"1230"`
	Difference     float64            `json:"difference" example:"25"`
	CountOperation int                `json:"count_operation" example:"6"`
	AVGIncome      float64            `json:"avg_income" example:"123.53"`
	AVGExpenditure float64            `json:"avg_expenditure" example:"632.12"`
	CostCategory   string             `json:"cost_category" example:"Medicine"`
	ShareCategory  map[string]float64 `json:"share_category" example:"Medicine:0.42,Medic:0.58"`
	MaxIncome      int                `json:"max_income" example:"1200"`
	MaxExpenditure int                `json:"max_expenditure" example:"155"`
}

// GetStatic			godoc
// @Summary 			Get static
// @Description 		Get all static about category and transaction
// @Tags 				static
// @Accept 				json
// @Produce 			json
// @Param				to				query int false "Time to"
// @Param				from 			query int false "Time from"
// @Param				category_id 	query int false "Category Id"
// @Success				200	{object}	StaticResponse "Get static successfully"
// @Failure 			400	{object}	response.ErrorResponse "Bad request"
// @Failure      		500 {object} response.ErrorResponse "Internal server error"
// @Router 				/static/{user_id}			[get]
func (h *StaticHTTPHandler) GetStatic(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	ResponseHandler := response.NewHandlerResponse(log, w)

	staticFilters, err := GetStaticFilters(r)
	if err != nil {
		ResponseHandler.ErrorResponse(err, "error getting static filters")
		return
	}

	staticDomain, err := h.StaticService.GetStatic(ctx, staticFilters)
	if err != nil {
		ResponseHandler.ErrorResponse(err, "error getting static domain")
		return
	}

	resp := StaticDomainToResponse(staticDomain)
	ResponseHandler.JSONResponseHandler(http.StatusOK, resp)
}

func StaticDomainToResponse(static core_domain.Static) StaticResponse {
	return StaticResponse{
		SumIncome:      static.SumIncome,
		SumExpenditure: static.SumExpenditure,
		Difference:     static.Difference,
		CountOperation: static.CountOperation,
		AVGIncome:      static.AVGIncome,
		AVGExpenditure: static.AVGExpenditure,
		CostCategory:   static.CostCategory,
		ShareCategory:  static.ShareCategory,
		MaxIncome:      static.MaxIncome,
		MaxExpenditure: static.MaxExpenditure,
	}
}

func GetStaticFilters(r *http.Request) (core_domain.FiltersStatic, error) {
	categoryID, err := core_http_utils.GetIntQueryParm(r, "category_id")
	if err != nil {
		return core_domain.FiltersStatic{}, fmt.Errorf("error while parsing category_id: %w", err)
	}
	to, err := core_http_utils.GetDateQueryParm(r, "to")
	if err != nil {
		return core_domain.FiltersStatic{}, fmt.Errorf("error getting query param: %w", err)
	}
	from, err := core_http_utils.GetDateQueryParm(r, "from")
	if err != nil {
		return core_domain.FiltersStatic{}, fmt.Errorf("error getting query param: %w", err)
	}
	userID, err := core_http_utils.GetValuePathInt(r, "user_id")
	if err != nil {
		return core_domain.FiltersStatic{}, fmt.Errorf("error getting query param: %w", err)
	}

	filters := core_domain.NewFiltersStatic(to, from, categoryID, userID)

	return filters, nil
}
