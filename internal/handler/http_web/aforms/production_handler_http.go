package aforms

import (
	"fgw_web_aforms/internal/handler"
	"fgw_web_aforms/internal/handler/http_err"
	"fgw_web_aforms/internal/handler/page"
	"fgw_web_aforms/internal/model"
	"fgw_web_aforms/internal/service"
	"fgw_web_aforms/pkg/common"
	"fgw_web_aforms/pkg/common/msg"
	"net/http"
)

const (
	tmplIndexHTML         = "index.html"
	tmplProductionHTML    = "productions.html"
	tmplProductionAddHTML = "web/html/aforms/production_add.html"
)

type ProductionHandlerHTML struct {
	productionService service.ProductionUserCase
	performerService  service.PerformerUseCase
	roleService       service.RoleUseCase
	logg              *common.Logger
	authMiddleware    *handler.AuthMiddleware
}

func NewProductionHandlerHTML(productionService service.ProductionUserCase, performerService service.PerformerUseCase,
	roleService service.RoleUseCase, logg *common.Logger, authMiddleware *handler.AuthMiddleware) *ProductionHandlerHTML {

	return &ProductionHandlerHTML{productionService, performerService, roleService, logg, authMiddleware}
}

func (p *ProductionHandlerHTML) ServeHTTPHTMLRouter(mux *http.ServeMux) {
	mux.HandleFunc("/aforms/productions", p.authMiddleware.RequireAuth(p.authMiddleware.RequireRole([]int{0, 4, 5}, p.AllProductionHTML)))
	mux.HandleFunc("/aforms/productions/add", p.authMiddleware.RequireAuth(p.authMiddleware.RequireRole([]int{0, 4, 5}, p.AddProductionHTML)))
}

func (p *ProductionHandlerHTML) AllProductionHTML(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if r.Method != http.MethodGet {
		http_err.SendErrorHTTP(w, http.StatusMethodNotAllowed, msg.H7000, p.logg, r)

		return
	}

	performerData, err := p.authMiddleware.GetUserData(r, p.performerService, p.roleService)
	if err != nil {
		http_err.SendErrorHTTP(w, http.StatusUnauthorized, msg.H7005, p.logg, r)

		return
	}

	productions, searchFields, sortFields, err := p.getProductions(w, r)
	if err != nil {
		http_err.SendErrorHTTP(w, http.StatusInternalServerError, msg.H7007, p.logg, r)

		return
	}

	data := page.NewDataPage("Варианты упаковки продукции", "productions", performerData,
		productions, sortFields, searchFields, true)

	page.RenderPages(w, tmplIndexHTML, data, r, tmplProductionHTML)
}

func (p *ProductionHandlerHTML) getProductions(w http.ResponseWriter, r *http.Request) ([]*model.Production, *page.SearchProductionsPage, *page.SortProductionsPage, error) {
	var productions []*model.Production
	var err error

	articlePattern := r.URL.Query().Get("articles")
	namePattern := r.URL.Query().Get("name")
	idPattern := r.URL.Query().Get("idProduction")

	sortField := r.URL.Query().Get("sort")
	sortOrder := r.URL.Query().Get("order")

	if articlePattern != "" || namePattern != "" || idPattern != "" {
		productions, err = p.productionService.SearchProductions(r.Context(), articlePattern, namePattern, idPattern)
		if err != nil {
			http_err.SendErrorHTTP(w, http.StatusNotFound, msg.H7008+err.Error(), p.logg, r)

			return nil, nil, nil, err
		}
	} else {
		productions, err = p.productionService.AllProductions(r.Context(), sortField, sortOrder)
		if err != nil {
			http_err.SendErrorHTTP(w, http.StatusNotFound, msg.H7000+err.Error(), p.logg, r)

			return nil, nil, nil, err
		}
	}

	return productions,
		&page.SearchProductionsPage{
			SearchArticle: articlePattern,
			SearchName:    namePattern,
			SearchId:      idPattern,
		},
		&page.SortProductionsPage{
			SortField: sortField,
			SortOrder: sortOrder,
		}, nil
}

func (p *ProductionHandlerHTML) AddProductionHTML(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	data := struct {
		Title string
	}{
		Title: "Добавить вариант упаковки",
	}

	page.RenderSinglePage(w, tmplProductionAddHTML, data, r)

}
