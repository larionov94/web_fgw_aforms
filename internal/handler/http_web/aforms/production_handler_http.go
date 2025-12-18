package aforms

import (
	"fgw_web_aforms/internal/handler"
	"fgw_web_aforms/internal/handler/http_err"
	"fgw_web_aforms/internal/service"
	"fgw_web_aforms/pkg/common"
	"net/http"
)

type ProductionHandlerHTML struct {
	productionService service.ProductionUserCase
	logg              *common.Logger
	authMiddleware    *handler.AuthMiddleware
}

func NewProductionHandlerHTML(productionService service.ProductionUserCase, logg *common.Logger, authMiddleware *handler.AuthMiddleware) *ProductionHandlerHTML {
	return &ProductionHandlerHTML{productionService, logg, authMiddleware}
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

}
