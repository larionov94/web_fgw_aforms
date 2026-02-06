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
	"strings"
)

const (
	tmplIndexHTML         = "index.html"
	tmplProductionHTML    = "productions.html"
	tmplProductionAddHTML = "production_add.html"
	tmplProductionUpdHTML = "production_upd.html"

	urlProductions = "/aforms/productions"

	renderPageTitle = "Варианты упаковки продукции"
	renderPageKey   = "productions"
	updPageTitle    = "Редактировать вариант упаковки"
	updPageKey      = "productionUpd"
	addPageTitle    = "Добавить вариант упаковки"
	addPageKey      = "productionAdd"
)

type ProductionHandlerHTML struct {
	productionService service.ProductionUseCase
	performerService  service.PerformerUseCase
	roleService       service.RoleUseCase
	catalogService    service.CatalogUseCase
	logg              *common.Logger
	authMiddleware    *handler.AuthMiddleware
	authPerformerData *http_web.AuthHandlerHTML
}

func NewProductionHandlerHTML(productionService service.ProductionUseCase, performerService service.PerformerUseCase,
	roleService service.RoleUseCase, catalogService service.CatalogUseCase, logg *common.Logger, authMiddleware *handler.AuthMiddleware,
	authPerformerData *http_web.AuthHandlerHTML) *ProductionHandlerHTML {

	return &ProductionHandlerHTML{productionService, performerService, roleService, catalogService, logg, authMiddleware, authPerformerData}
}

func (p *ProductionHandlerHTML) ServeHTTPHTMLRouter(mux *http.ServeMux) {
	mux.HandleFunc("/aforms/productions", p.authMiddleware.RequireAuth(p.authMiddleware.RequireRole([]int{0, 4, 5}, p.RenderProductionsPage)))
	mux.HandleFunc("/aforms/productions/add", p.authMiddleware.RequireAuth(p.authMiddleware.RequireRole([]int{0, 4, 5}, p.AddProductionForm)))
	mux.HandleFunc("/aforms/productions/upd", p.authMiddleware.RequireAuth(p.authMiddleware.RequireRole([]int{0, 4, 5}, p.UpdateProductionForm)))
}

// RenderProductionsPage отображает страницу продукции.
func (p *ProductionHandlerHTML) RenderProductionsPage(w http.ResponseWriter, r *http.Request) {
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

	productions, searchFields, sortFields, err := p.fetchProductionsWithParams(w, r)
	if err != nil {
		http_err.SendErrorHTTP(w, http.StatusInternalServerError, msg.H7007, p.logg, r)

		return
	}

	data := page.NewDataPage(
		renderPageTitle,
		renderPageKey,
		performerData,
		productions,
		sortFields,
		searchFields,
		true,
		nil,
		nil,
		nil,
		nil,
		nil,
	)

	page.RenderPages(w, tmplIndexHTML, data, r, tmplProductionHTML, tmplProductionAddHTML, tmplProductionUpdHTML, tmplPlanHTML)
}

// fetchProductionsWithParams - получить продукцию с учетом параметров запроса.
func (p *ProductionHandlerHTML) fetchProductionsWithParams(w http.ResponseWriter, r *http.Request) ([]*model.Production,
	*page.SearchProductionsPage, *page.SortProductionsPage, error) {

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

func (p *ProductionHandlerHTML) UpdateProductionForm(w http.ResponseWriter, r *http.Request) {
	performerData, err := p.authPerformerData.AuthenticatePerformer(r)
	if err != nil {
		http_err.SendErrorHTTP(w, http.StatusUnauthorized, msg.H7005, p.logg, r)

		return
	}

	switch r.Method {
	case http.MethodGet:
		p.handleGetUpdateForm(w, r, performerData)
	case http.MethodPost:
		p.handlePostUpdateForm(w, r, performerData)
	default:
		http.Error(w, msg.H7000, http.StatusMethodNotAllowed)

		return
	}
}

// handleGetUpdateForm - обработчик формы для GET запроса обновления продукции.
func (p *ProductionHandlerHTML) handleGetUpdateForm(w http.ResponseWriter, r *http.Request, performerData *handler.PerformerData) {
	production, err := p.ensureProductionExists(r, "idProduction")
	if err != nil {
		http_err.SendErrorHTTP(w, http.StatusNotFound, msg.H7000+err.Error(), p.logg, r)

		return
	}

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
		updPageTitle,
		updPageKey,
		performerData,
		[]*model.Production{production},
		nil,
		nil,
		false,
		designNameList,
		colorList,
		nil,
		nil,
		nil,
	)

	page.RenderPages(w, tmplIndexHTML, data, r, tmplPlanHTML, tmplProductionHTML, tmplProductionAddHTML, tmplProductionUpdHTML)

	return
}

// handlePostUpdateForm - обработчик формы для POST запроса обновления продукции.
func (p *ProductionHandlerHTML) handlePostUpdateForm(w http.ResponseWriter, r *http.Request, performerData *handler.PerformerData) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if err := r.ParseForm(); err != nil {
		http_err.SendErrorHTTP(w, http.StatusBadRequest, msg.H7018+err.Error(), p.logg, r)

		return
	}

	idProduction, err := p.extractIdParamFormValue(r, "idProduction")
	if err != nil {
		http_err.SendErrorHTTP(w, http.StatusBadRequest, msg.H7104, p.logg, r)

		return
	}

	prPartLastDate := strings.TrimSpace(r.FormValue("PrPartLastDate"))
	formatPrPartLastDate, err := convert.ParseToMSSQLDateTime(prPartLastDate)
	if err != nil {
		page.RenderErrorPage(w, 400, msg.H7101, r)

		return
	}

	product := &model.Production{
		PrName:         strings.TrimSpace(r.FormValue("PrName")),
		PrShortName:    strings.TrimSpace(r.FormValue("PrShortName")),
		PrPackName:     strings.TrimSpace(r.FormValue("PrPackName")),
		PrArticle:      strings.TrimSpace(r.FormValue("PrArticle")),
		PrColor:        strings.TrimSpace(r.FormValue("PrColor")),
		PrCount:        convert.ParseFormFieldInt(r, "PrCount"),
		PrRows:         convert.ParseFormFieldInt(r, "PrRows"),
		PrWeight:       convert.ParseFormFieldFloat(r, "PrWeight"),
		PrHWD:          strings.TrimSpace(r.FormValue("PrHWD")),
		PrInfo:         strings.TrimSpace(r.FormValue("PrInfo")),
		PrPart:         convert.ParseFormFieldInt(r, "PrPart"),
		PrPartLastDate: formatPrPartLastDate, // Дата выпуска продукции, дата идет на этикетку.
		PrPartAutoInc:  convert.ParseFormFieldInt(r, "PrPartAutoInc"),
		PrPerGodn:      convert.ParseFormFieldInt(r, "PrPerGodn"),
		PrSAP:          strings.TrimSpace(r.FormValue("PrSAP")),
		PrProdType:     convert.ParseFormFieldBool(r, "PrProdType"),
		PrUmbrella:     convert.ParseFormFieldBool(r, "PrUmbrella"),
		PrPerfumery:    convert.ParseFormFieldBool(r, "PrPerfumery"),
		PrSun:          convert.ParseFormFieldBool(r, "PrSun"),
		PrDecl:         convert.ParseFormFieldBool(r, "PrDecl"),
		PrParty:        convert.ParseFormFieldBool(r, "PrParty"),
		PrGL:           convert.ParseFormFieldInt(r, "PrGL"),
		AuditRec: model.Audit{
			CreatedBy: performerData.PerformerId,
			UpdatedBy: performerData.PerformerId,
		},
	}

	if err = p.productionService.UpdProduction(r.Context(), idProduction, product); err != nil {
		http_err.SendErrorHTTP(w, http.StatusInternalServerError, msg.H7000+err.Error(), p.logg, r)

		return
	}
	http.Redirect(w, r, urlProductions, http.StatusSeeOther)

	return
}

// extractIdParam - извлечь параметр идентификатора.
func (p *ProductionHandlerHTML) extractIdParam(r *http.Request, paramId string) (int, error) {
	idParamStr := r.URL.Query().Get(paramId)
	if idParamStr == "" {
		return 0, fmt.Errorf("%s idParamStr=%s", msg.H7102, idParamStr)
	}

	idParam := convert.ConvStrToInt(idParamStr)
	if idParam == 0 {
		return 0, fmt.Errorf("%s idParam=%d", msg.H7104, idParam)
	}

	return idParam, nil
}

// extractIdParamFormValue - извлечь параметр идентификатора из формы.
func (p *ProductionHandlerHTML) extractIdParamFormValue(r *http.Request, paramId string) (int, error) {
	idParamStr := r.FormValue(paramId)
	if idParamStr == "" {
		return 0, fmt.Errorf("%s idParamStr=%s", msg.H7102, idParamStr)
	}

	idParam := convert.ConvStrToInt(idParamStr)
	if idParam == 0 {
		return 0, fmt.Errorf("%s idParam=%d", msg.H7104, idParam)
	}

	return idParam, nil
}

// ensureProductionExists - обеспечить существование продукции.
func (p *ProductionHandlerHTML) ensureProductionExists(r *http.Request, paramId string) (*model.Production, error) {
	idProduction, err := p.extractIdParam(r, paramId)
	if err != nil {
		return nil, fmt.Errorf("%s idProduction=%d", msg.H7102, idProduction)
	}

	production, err := p.productionService.FindByIdProduction(r.Context(), idProduction)
	if err != nil {
		return nil, fmt.Errorf("%s: %s", msg.H7103, err.Error())
	}

	if production == nil {
		return nil, fmt.Errorf("%s", msg.H7103)
	}

	return production, nil
}

func (p *ProductionHandlerHTML) AddProductionForm(w http.ResponseWriter, r *http.Request) {
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

// handlerGetAddForm - обработчик формы для GET запроса добавления продукции.
func (p *ProductionHandlerHTML) handlerGetAddForm(w http.ResponseWriter, r *http.Request, performerData *handler.PerformerData) {
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

	if r.URL.Query().Get("basedOn") != "" {
		basedOnProduction, err := p.ensureProductionExists(r, "basedOn")
		if err != nil {
			http_err.SendErrorHTTP(w, http.StatusNotFound, msg.H7000+err.Error(), p.logg, r)

			return
		}

		basedOnProduction.IdProduction = 0
		basedOnProduction.PrArticle = ""

		data := page.NewDataPage(
			addPageTitle,
			addPageKey,
			performerData,
			[]*model.Production{basedOnProduction},
			nil,
			nil,
			false,
			designNameList,
			colorList,
			nil,
			nil,
			nil,
		)

		page.RenderPages(w, tmplIndexHTML, data, r, tmplPlanHTML, tmplProductionHTML, tmplProductionAddHTML, tmplProductionUpdHTML)

		return

	}

	data := page.NewDataPage(
		addPageTitle,
		addPageKey,
		performerData,
		nil,
		nil,
		nil,
		false,
		designNameList,
		colorList,
		nil,
		nil,
		nil,
	)

	page.RenderPages(w, tmplIndexHTML, data, r, tmplPlanHTML, tmplProductionHTML, tmplProductionAddHTML, tmplProductionUpdHTML)

	return
}

// AddProductionForm - форма для добавления продукции.
func (p *ProductionHandlerHTML) handlerPostAddForm(w http.ResponseWriter, r *http.Request, performerData *handler.PerformerData) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if err := r.ParseForm(); err != nil {
		http_err.SendErrorHTTP(w, http.StatusBadRequest, msg.H7018+err.Error(), p.logg, r)

		return
	}

	prArticle := strings.TrimSpace(r.FormValue("PrArticle"))
	if len(prArticle) != 2 {
		page.RenderErrorPage(w, 400, msg.H7100, r)

		return
	}

	prPartLastDate := strings.TrimSpace(r.FormValue("PrPartLastDate"))
	formatPrPartLastDate, err := convert.ParseToMSSQLDateTime(prPartLastDate)
	if err != nil {
		page.RenderErrorPage(w, 400, msg.H7101, r)

		return
	}

	product := &model.Production{
		PrName:         strings.TrimSpace(r.FormValue("PrName")),
		PrShortName:    strings.TrimSpace(r.FormValue("PrShortName")),
		PrPackName:     strings.TrimSpace(r.FormValue("PrPackName")),
		PrArticle:      prArticle,
		PrColor:        strings.TrimSpace(r.FormValue("PrColor")),
		PrCount:        convert.ParseFormFieldInt(r, "PrCount"),
		PrRows:         convert.ParseFormFieldInt(r, "PrRows"),
		PrWeight:       convert.ParseFormFieldFloat(r, "PrWeight"),
		PrHWD:          strings.TrimSpace(r.FormValue("PrHWD")),
		PrInfo:         strings.TrimSpace(r.FormValue("PrInfo")),
		PrPart:         convert.ParseFormFieldInt(r, "PrPart"),
		PrPartLastDate: formatPrPartLastDate, // Дата выпуска продукции, дата идет на этикетку.
		PrPartAutoInc:  convert.ParseFormFieldInt(r, "PrPartAutoInc"),
		PrPerGodn:      convert.ParseFormFieldInt(r, "PrPerGodn"),
		PrSAP:          strings.TrimSpace(r.FormValue("PrSAP")),
		PrProdType:     convert.ParseFormFieldBool(r, "PrProdType"),
		PrUmbrella:     convert.ParseFormFieldBool(r, "PrUmbrella"),
		PrPerfumery:    convert.ParseFormFieldBool(r, "PrPerfumery"),
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
