package aforms

import (
	"fgw_web_aforms/internal/handler"
	"fgw_web_aforms/internal/handler/http_err"
	"fgw_web_aforms/internal/handler/http_web"
	"fgw_web_aforms/internal/handler/page"
	"fgw_web_aforms/internal/model"
	"fgw_web_aforms/internal/service"
	"fgw_web_aforms/pkg/common"
	"fgw_web_aforms/pkg/common/msg"
	"net/http"
	"time"
)

const (
	tmplPlanHTML = "plans.html"

	renderPagePlanTitle = "Планы сменно-суточных заданий"
	renderPagePlanKey   = "plans"
)

type PlanHandlerHTML struct {
	planService       service.PlanUseCase
	logg              *common.Logger
	authMiddleware    *handler.AuthMiddleware
	authPerformerData *http_web.AuthHandlerHTML
}

func NewPlanHandlerHTML(planService service.PlanUseCase, logg *common.Logger, authMiddleware *handler.AuthMiddleware,
	authPerformerData *http_web.AuthHandlerHTML) *PlanHandlerHTML {
	return &PlanHandlerHTML{planService: planService, logg: logg, authMiddleware: authMiddleware, authPerformerData: authPerformerData}
}

func (p *PlanHandlerHTML) ServeHTTPHTMLRouter(mux *http.ServeMux) {
	mux.HandleFunc("/aforms/plans", p.authMiddleware.RequireAuth(p.authMiddleware.RequireRole([]int{0, 4, 5}, p.RenderPlanPage)))
}

func (p *PlanHandlerHTML) RenderPlanPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if r.Method != http.MethodGet {
		http_err.SendErrorHTTP(w, http.StatusMethodNotAllowed, msg.H7000, p.logg, r)

		return
	}

	performerData, err := p.authPerformerData.AuthenticatePerformer(r)
	if err != nil {
		http_err.SendErrorHTTP(w, http.StatusUnauthorized, msg.H7005, p.logg, r)

		return
	}

	plans, sortField, err := p.fetchPlansWithParams(w, r)
	if err != nil {
		http_err.SendErrorHTTP(w, http.StatusInternalServerError, msg.H7007, p.logg, r)

		return
	}

	data := page.NewDataPage(
		renderPagePlanTitle,
		renderPagePlanKey,
		performerData,
		nil,
		nil,
		nil,
		false,
		nil,
		nil,
		sortField,
		plans,
	)

	page.RenderPages(w, tmplIndexHTML, data, r, tmplPlanHTML, tmplProductionHTML, tmplProductionAddHTML, tmplProductionUpdHTML)
}

// fetchPlansWithParams - получить план с учетом параметров запроса.
func (p *PlanHandlerHTML) fetchPlansWithParams(w http.ResponseWriter, r *http.Request) ([]*model.Plan,
	*page.SortPlanPage, error) {

	var plans []*model.Plan
	var err error

	sortField := r.URL.Query().Get("sort")
	sortOrder := r.URL.Query().Get("order")

	startDate := r.URL.Query().Get("startDate")
	endDate := r.URL.Query().Get("endDate")

	if startDate == "" {
		startDate = time.Now().AddDate(0, 0, -180).Format("2006-01-02")
	}

	if endDate == "" {
		endDate = time.Now().Format("2006-01-02")
	}

	plans, err = p.planService.AllPlans(r.Context(), sortField, sortOrder, startDate, endDate)
	if err != nil {
		http_err.SendErrorHTTP(w, http.StatusNotFound, msg.H7000+err.Error(), p.logg, r)
		return nil, nil, err
	}

	return plans,
		&page.SortPlanPage{
			SortField: sortField,
			SortOrder: sortOrder,
			StartDate: startDate,
			EndDate:   endDate,
		}, nil
}
