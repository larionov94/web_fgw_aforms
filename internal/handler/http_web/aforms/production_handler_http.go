package aforms

import (
	"fgw_web_aforms/internal/handler"
	"fgw_web_aforms/internal/handler/http_err"
	"fgw_web_aforms/internal/handler/http_web"
	"fgw_web_aforms/internal/handler/page"
	"fgw_web_aforms/internal/service"
	"fgw_web_aforms/pkg/common"
	"net/http"
)

const (
	tmplIndexHTML      = "index.html"
	tmplProductionHTML = "productions.html"
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
}

func (p *ProductionHandlerHTML) AllProductionHTML(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if r.Method != http.MethodGet {
		http_err.SendErrorHTTP(w, http.StatusMethodNotAllowed, "", p.logg, r)

		return
	}

	performerFIO, performerId, roleName, err := p.authMiddleware.GetUserData(r, p.performerService, p.roleService)
	if err != nil {

		return
	}

	data := http_web.NewDataPage("Варианты упаковки продукции", "productions", performerFIO, performerId, roleName)

	page.RenderPages(w, tmplIndexHTML, data, r, tmplProductionHTML)
}
