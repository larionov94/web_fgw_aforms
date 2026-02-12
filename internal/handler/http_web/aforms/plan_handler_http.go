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
	"fgw_web_aforms/pkg/convert"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	tmplPlanHTML    = "plans.html"
	tmplPlanAddHTML = "plan_add.html"

	urlPlan = "/aforms/plans"

	renderPagePlanTitle = "Планы сменно-суточных заданий"
	renderPagePlanKey   = "plans"

	addPlanPageTitle  = "Добавить сменно-суточный план"
	addPlanPageAddKey = "planAdd"
)

type PlanHandlerHTML struct {
	planService       service.PlanUseCase
	logg              *common.Logger
	authMiddleware    *handler.AuthMiddleware
	authPerformerData *http_web.AuthHandlerHTML
	productionService service.ProductionUseCase
	sectorService     service.SectorUseCase
}

func NewPlanHandlerHTML(planService service.PlanUseCase, logg *common.Logger, authMiddleware *handler.AuthMiddleware,
	authPerformerData *http_web.AuthHandlerHTML, productionService service.ProductionUseCase, sectorService service.SectorUseCase) *PlanHandlerHTML {
	return &PlanHandlerHTML{planService: planService, logg: logg, authMiddleware: authMiddleware, authPerformerData: authPerformerData, productionService: productionService, sectorService: sectorService}
}

func (p *PlanHandlerHTML) ServeHTTPHTMLRouter(mux *http.ServeMux) {
	mux.HandleFunc("/aforms/plans", p.authMiddleware.RequireAuth(p.authMiddleware.RequireRole([]int{0, 4, 5}, p.RenderPlanPage)))
	mux.HandleFunc("/aforms/plans/add", p.authMiddleware.RequireAuth(p.authMiddleware.RequireRole([]int{0, 4, 5}, p.AddPlanForm)))
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

	productions, err := p.productionService.AllProductions(r.Context(), "", "")
	if err != nil {
		http_err.SendErrorHTTP(w, http.StatusNotFound, msg.H7000+err.Error(), p.logg, r)

		return
	}

	sectors, err := p.sectorService.AllSector(r.Context())
	if err != nil {
		http_err.SendErrorHTTP(w, http.StatusNotFound, msg.H7000+err.Error(), p.logg, r)

		return
	}

	data := page.NewDataPage(
		renderPagePlanTitle,
		renderPagePlanKey,
		performerData,
		productions,
		nil,
		nil,
		false,
		nil,
		nil,
		sortField,
		plans,
		sectors,
	)

	page.RenderPages(w, tmplIndexHTML, data, r, tmplPlanHTML, tmplProductionHTML, tmplProductionAddHTML, tmplProductionUpdHTML, tmplPlanAddHTML)
}

// fetchPlansWithParams - получить план с учетом параметров запроса.
func (p *PlanHandlerHTML) fetchPlansWithParams(w http.ResponseWriter, r *http.Request) ([]*model.Plan,
	*page.SortPlanPage, error) {

	var plans []*model.Plan
	var err error
	var idProduction *int
	var idSector *int

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

	idProductionStr := r.URL.Query().Get("idProduction")
	if idProductionStr != "" {
		idProd, err := strconv.Atoi(idProductionStr)
		if err != nil {
			p.logg.LogW(fmt.Sprintf("Некорректный idProduction: %s", idProductionStr))
			idProduction = nil
		} else {
			idProduction = &idProd
		}
	}

	idSectorStr := r.URL.Query().Get("idSector")
	if idSectorStr != "" {
		idSec, err := strconv.Atoi(idSectorStr)
		if err != nil {
			p.logg.LogW(fmt.Sprintf("Некорректный idSector: %s", idSectorStr))
			idSector = nil
		} else {
			idSector = &idSec
		}
	}

	prName := r.URL.Query().Get("PrName")
	secName := r.URL.Query().Get("SectorName")

	plans, err = p.planService.AllPlans(r.Context(), sortField, sortOrder, startDate, endDate, idProduction, idSector)
	if err != nil {
		http_err.SendErrorHTTP(w, http.StatusNotFound, msg.H7000+err.Error(), p.logg, r)
		return nil, nil, err
	}

	return plans,
		&page.SortPlanPage{
			SortField:    sortField,
			SortOrder:    sortOrder,
			StartDate:    startDate,
			EndDate:      endDate,
			IdProduction: idProduction,
			IdSector:     idSector,
			PrName:       prName,
			SectorName:   secName,
		}, nil
}

func (p *PlanHandlerHTML) AddPlanForm(w http.ResponseWriter, r *http.Request) {
	performerData, err := p.authPerformerData.AuthenticatePerformer(r)
	if err != nil {
		http_err.SendErrorHTTP(w, http.StatusUnauthorized, msg.H7005, p.logg, r)

		return
	}

	switch r.Method {
	case http.MethodGet:
		p.handlerGetAddForm(w, r, performerData)
	case http.MethodPost:
		p.handlerPostAddForm(w, r, performerData)
	default:
		http.Error(w, msg.H7000, http.StatusMethodNotAllowed)

		return
	}
}

func (p *PlanHandlerHTML) handlerPostAddForm(w http.ResponseWriter, r *http.Request, performerData *handler.PerformerData) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := r.ParseForm(); err != nil {
		http_err.SendErrorHTTP(w, http.StatusBadRequest, msg.H7018+err.Error(), p.logg, r)

		return
	}

	planDate := strings.TrimSpace(r.FormValue("PlanDate"))

	formatPlanDate, err := convert.ParseToMSSQLDateTime(planDate)
	if err != nil {
		page.RenderErrorPage(w, 400, msg.H7101, r)

		return
	}

	plan := &model.Plan{
		PlanShift:     convert.ParseFormFieldInt(r, "PlanShift"),
		ExtProduction: convert.ParseFormFieldInt(r, "extProduction"),
		ExtSector:     convert.ParseFormFieldInt(r, "extSector"),
		PlanCount:     convert.ParseFormFieldInt(r, "PlanCount"),
		PlanDate:      formatPlanDate,
		PlanInfo:      r.FormValue("PlanInfo"),
		AuditRec: model.Audit{
			CreatedBy: performerData.PerformerId,
			UpdatedBy: performerData.PerformerId,
		},
	}

	if err := p.planService.AddPlan(r.Context(), plan); err != nil {
		http_err.SendErrorHTTP(w, http.StatusInternalServerError, msg.H7000+err.Error(), p.logg, r)

		return
	}
	http.Redirect(w, r, urlPlan, http.StatusSeeOther)

	return
}

func (p *PlanHandlerHTML) handlerGetAddForm(w http.ResponseWriter, r *http.Request, performerData *handler.PerformerData) {
	productions, err := p.productionService.AllProductions(r.Context(), "", "")
	if err != nil {
		http_err.SendErrorHTTP(w, http.StatusNotFound, msg.H7000+err.Error(), p.logg, r)

		return
	}

	sectors, err := p.sectorService.AllSector(r.Context())
	if err != nil {
		http_err.SendErrorHTTP(w, http.StatusNotFound, msg.H7000+err.Error(), p.logg, r)

		return
	}

	data := page.NewDataPage(
		addPlanPageTitle,
		addPlanPageAddKey,
		performerData,
		productions,
		nil,
		nil,
		false,
		nil,
		nil,
		nil,
		nil,
		sectors,
	)

	page.RenderPages(w, tmplIndexHTML, data, r, tmplPlanHTML, tmplProductionAddHTML, tmplProductionUpdHTML, tmplPlanAddHTML, tmplProductionHTML)

	return
}
