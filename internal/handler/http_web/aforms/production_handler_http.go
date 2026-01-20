package aforms

import (
	"fgw_web_aforms/internal/handler"
	"fgw_web_aforms/internal/handler/http_err"
	"fgw_web_aforms/internal/handler/page"
	"fgw_web_aforms/internal/model"
	"fgw_web_aforms/internal/service"
	"fgw_web_aforms/pkg/common"
	"fgw_web_aforms/pkg/common/msg"
	"fgw_web_aforms/pkg/convert"
	"net/http"
	"strings"
)

const (
	tmplIndexHTML         = "index.html"
	tmplProductionHTML    = "productions.html"
	tmplProductionAddHTML = "production_add.html"

	urlProductions = "/aforms/productions"
)

type ProductionHandlerHTML struct {
	productionService service.ProductionUseCase
	performerService  service.PerformerUseCase
	roleService       service.RoleUseCase
	catalogService    service.CatalogUseCase
	logg              *common.Logger
	authMiddleware    *handler.AuthMiddleware
}

func NewProductionHandlerHTML(productionService service.ProductionUseCase, performerService service.PerformerUseCase,
	roleService service.RoleUseCase, catalogService service.CatalogUseCase, logg *common.Logger, authMiddleware *handler.AuthMiddleware) *ProductionHandlerHTML {

	return &ProductionHandlerHTML{productionService, performerService, roleService, catalogService, logg, authMiddleware}
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
		productions, sortFields, searchFields, true, nil, nil)

	page.RenderPages(w, tmplIndexHTML, data, r, tmplProductionHTML, tmplProductionAddHTML)
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

	performerData, err := p.authMiddleware.GetUserData(r, p.performerService, p.roleService)
	if err != nil {
		http_err.SendErrorHTTP(w, http.StatusUnauthorized, msg.H7005, p.logg, r)

		return
	}

	// Обработка GET запроса - отображение формы
	if r.Method == http.MethodGet {
		designNameList, err := p.catalogService.DesignNameAll(r.Context())
		if err != nil {
			http_err.SendErrorHTTP(w, http.StatusInternalServerError, msg.H7000+err.Error(), p.logg, r)
			return
		}

		colorList, err := p.catalogService.ColorAll(r.Context())
		if err != nil {
			http_err.SendErrorHTTP(w, http.StatusInternalServerError, msg.H7000+err.Error(), p.logg, r)
			return
		}

		data := page.NewDataPage(
			"Добавить вариант упаковки",
			"productionAdd",
			performerData,
			nil,
			nil,
			nil,
			false,
			designNameList,
			colorList)

		page.RenderPages(w, tmplIndexHTML, data, r, tmplProductionHTML, tmplProductionAddHTML)

		return
	}

	// Обработка POST запроса - сохранение данных
	if r.Method == http.MethodPost {
		if err := r.ParseForm(); err != nil {
			http_err.SendErrorHTTP(w, http.StatusBadRequest, msg.H7018+err.Error(), p.logg, r)

			return
		}

		prArticle := strings.TrimSpace(r.FormValue("PrArticle"))
		if len(prArticle) != 2 {
			page.RenderErrorPage(w, 400, msg.H7100, r)

			return
		}

		// Создаем продукт из данных формы
		product := &model.Production{
			PrName:         strings.TrimSpace(r.FormValue("PrName")),
			PrShortName:    strings.TrimSpace(r.FormValue("PrShortName")),
			PrPackName:     strings.TrimSpace(r.FormValue("PrPackName")),
			PrType:         "",
			PrArticle:      prArticle,
			PrColor:        strings.TrimSpace(r.FormValue("PrColor")),
			PrBarCode:      "",
			PrCount:        convert.ParseFormFieldInt(r, "PrCount"),
			PrRows:         convert.ParseFormFieldInt(r, "PrRows"),
			PrWeight:       convert.ParseFormFieldFloat(r, "PrWeight"),
			PrHWD:          strings.TrimSpace(r.FormValue("PrHWD")),
			PrInfo:         "",
			PrStatus:       "",
			PrEditDate:     "",
			PrEditUser:     0,
			PrPart:         0,
			PrPartLastDate: "",
			PrPartAutoInc:  0,
			PrPartRealDate: "",
			PrArchive:      false,
			PrPerGodn:      0,
			PrSAP:          strings.TrimSpace(r.FormValue("PrSAP")),
			PrProdType:     false,
			PrUmbrella:     convert.ParseFormFieldBool(r, "PrUmbrella"),
			PrSun:          convert.ParseFormFieldBool(r, "PrSun"),
			PrDecl:         convert.ParseFormFieldBool(r, "PrDecl"),
			PrParty:        convert.ParseFormFieldBool(r, "PrParty"),
			PrGL:           convert.ParseFormFieldInt(r, "PrGL"),
			AuditRec: model.Audit{
				CreatedBy: performerData.PerformerId,
				UpdatedBy: performerData.PerformerId,
			},
		}

		if err := p.productionService.AddProduction(r.Context(), product); err != nil {
			http_err.SendErrorHTTP(w, http.StatusInternalServerError, msg.H7000+err.Error(), p.logg, r)

			return
		}
		http.Redirect(w, r, urlProductions, http.StatusSeeOther)

		return
	}

	http.Error(w, msg.H7000, http.StatusMethodNotAllowed)
}
