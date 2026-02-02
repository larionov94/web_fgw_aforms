package page

import (
	"fgw_web_aforms/internal/handler"
	"fgw_web_aforms/internal/model"
	"fgw_web_aforms/pkg/common/msg"
	"fgw_web_aforms/pkg/convert"
	"fmt"
	"html/template"
	"net/http"
)

const (
	prefixDefaultTmpl = "web/html/"
	prefixAFormsTmpl  = "web/html/aforms/"
	tmplErrorHTML     = "error.html"
)

type SortProductionsPage struct {
	SortField string
	SortOrder string
}

type SortPlanPage struct {
	SortField string
	SortOrder string
	StartDate string
	EndDate   string
}

type SearchProductionsPage struct {
	SearchArticle string
	SearchName    string
	SearchId      string
}

type DataPage struct {
	Title          string                 // Title - название страницы
	CurrentPage    string                 // CurrentPage - шаблон страницы в html
	InfoPerformer  *handler.PerformerData // InfoPerformer - информация об авторизованном сотруднике
	Productions    []*model.Production    // Productions - список продукции
	SortProducts   *SortProductionsPage   // SortProducts - сортировка продукции
	SearchProducts *SearchProductionsPage // SearchProducts - фильтр продукции
	IsSearch       bool                   // IsSearch - разрешить поиск
	DesignNameList []*model.Catalog       // DesignNameList - список конструкторских наименований
	ColorList      []*model.Catalog       // ColorList - список цветов продукции
	SortPlans      *SortPlanPage          // SortPlans - сортировка плана
	Plans          []*model.Plan          // Plans - список планов
}

func NewDataPage(
	title string,
	currentPage string,
	infoPerformer *handler.PerformerData,
	productions []*model.Production,
	sortProducts *SortProductionsPage,
	searchProductions *SearchProductionsPage,
	isSearch bool,
	designNameList, colorList []*model.Catalog,
	sortPlan *SortPlanPage,
	plans []*model.Plan) *DataPage {

	return &DataPage{
		title,
		currentPage,
		infoPerformer,
		productions,
		sortProducts,
		searchProductions,
		isSearch,
		designNameList,
		colorList,
		sortPlan,
		plans}
}
func SetSecureHTMLHeaders(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, private, max-age=0")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("X-Frame-Options", "DENY")
}

// Вспомогательная функция для рендеринга ошибки напрямую
func renderErrorDirectly(w http.ResponseWriter, statusCode int, msgCode string, r *http.Request) {
	SetSecureHTMLHeaders(w)
	w.WriteHeader(statusCode)

	// Простой HTML для ошибки (без рекурсии)
	errorHTML := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <title>Ошибка %d</title>
    <meta charset="utf-8">
    <style>
        body { font-family: Arial, sans-serif; padding: 20px; }
        .error { color: #d9534f; border: 1px solid #d9534f; padding: 15px; margin: 10px 0; }
    </style>
</head>
<body>
    <h1>Ошибка %d</h1>
    <div class="error">
        <strong>Сообщение:</strong> %s<br>
        <strong>Метод:</strong> %s<br>
        <strong>Путь:</strong> %s
    </div>
</body>
</html>`, statusCode, statusCode, msgCode, r.Method, r.URL.Path)

	_, err := fmt.Fprint(w, errorHTML)
	if err != nil {
		return
	}
}

func RenderPage(w http.ResponseWriter, tmpl string, data interface{}, r *http.Request) {
	templatePath := prefixDefaultTmpl + tmpl

	parseTmpl, err := template.New(tmpl).Funcs(
		template.FuncMap{
			"formatDateTime": convert.FormatDateTime,
			"formatDate":     convert.FormatDate,
			"buildSortURL":   convert.BuildSortURL,
		}).ParseFiles(templatePath)

	if err != nil {
		renderErrorDirectly(w, http.StatusInternalServerError, msg.H7002+err.Error(), r)

		return
	}

	if err = parseTmpl.ExecuteTemplate(w, tmpl, data); err != nil {
		renderErrorDirectly(w, http.StatusInternalServerError, msg.H7003+err.Error(), r)

		return
	}
}

func RenderPages(w http.ResponseWriter, tmpl string, data interface{}, r *http.Request, additionalTemplates ...string) {
	templatePaths := []string{prefixDefaultTmpl + tmpl}

	for _, additionalTmpl := range additionalTemplates {
		templatePaths = append(templatePaths, prefixAFormsTmpl+additionalTmpl)
	}

	parseTmpl, err := template.New(tmpl).Funcs(
		template.FuncMap{
			"formatDateTime":      convert.FormatDateTime,
			"splitDimensions":     convert.SplitDimensions,
			"formatDateTimeLocal": convert.FormatDateTimeLocal,
			"formatDate":          convert.FormatDate,
			"buildSortURL":        convert.BuildSortURL,
		}).ParseFiles(templatePaths...)

	if err != nil {
		renderErrorDirectly(w, http.StatusInternalServerError, msg.H7002+err.Error(), r)

		return
	}

	if err = parseTmpl.ExecuteTemplate(w, tmpl, data); err != nil {
		renderErrorDirectly(w, http.StatusInternalServerError, msg.H7003+err.Error(), r)

		return
	}
}

func RenderErrorPage(w http.ResponseWriter, statusCode int, msgCode string, r *http.Request) {
	SetSecureHTMLHeaders(w)

	templatePath := prefixDefaultTmpl + tmplErrorHTML

	parseTmpl, err := template.New(tmplErrorHTML).Funcs(
		template.FuncMap{
			"formatDateTime": convert.FormatDateTime,
			"formatDate":     convert.FormatDate,
			"buildSortURL":   convert.BuildSortURL,
		}).ParseFiles(templatePath)

	if err != nil {
		renderErrorDirectly(w, statusCode, msgCode, r)

		return
	}

	data := struct {
		Title      string
		MsgCode    string
		StatusCode int
		Method     string
		Path       string
	}{
		Title:      "Ошибка",
		MsgCode:    msgCode,
		StatusCode: statusCode,
		Method:     r.Method,
		Path:       r.URL.Path,
	}

	if err = parseTmpl.Execute(w, data); err != nil {
		renderErrorDirectly(w, statusCode, msgCode, r)
		return
	}
}
