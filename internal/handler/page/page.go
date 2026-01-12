package page

import (
	"fgw_web_aforms/internal/handler"
	"fgw_web_aforms/internal/model"
	"fgw_web_aforms/pkg/common/msg"
	"fgw_web_aforms/pkg/convert"
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

type SearchProductionsPage struct {
	SearchArticle string
	SearchName    string
	SearchId      string
}

type DataPage struct {
	Title          string
	CurrentPage    string
	InfoPerformer  *handler.PerformerData
	Productions    []*model.Production
	SortProducts   *SortProductionsPage
	SearchProducts *SearchProductionsPage
	IsSearch       bool
}

func NewDataPage(title string, currentPage string, infoPerformer *handler.PerformerData, productions []*model.Production,
	sortProducts *SortProductionsPage, searchProductions *SearchProductionsPage, isSearch bool) *DataPage {

	return &DataPage{title, currentPage, infoPerformer, productions, sortProducts, searchProductions, isSearch}
}
func SetSecureHTMLHeaders(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, private, max-age=0")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("X-Frame-Options", "DENY")
}

func RenderPage(w http.ResponseWriter, tmpl string, data interface{}, r *http.Request) {
	templatePath := prefixDefaultTmpl + tmpl

	parseTmpl, err := template.New(tmpl).Funcs(
		template.FuncMap{
			"formatDateTime": convert.FormatDateTime,
		}).ParseFiles(templatePath)

	if err != nil {
		RenderErrorPage(w, http.StatusInternalServerError, msg.H7002+err.Error(), r)

		return
	}

	if err = parseTmpl.ExecuteTemplate(w, tmpl, data); err != nil {
		RenderErrorPage(w, http.StatusInternalServerError, msg.H7003+err.Error(), r)

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
			"formatDateTime": convert.FormatDateTime,
		}).ParseFiles(templatePaths...)

	if err != nil {
		RenderErrorPage(w, http.StatusInternalServerError, msg.H7002+err.Error(), r)

		return
	}

	if err = parseTmpl.ExecuteTemplate(w, tmpl, data); err != nil {
		RenderErrorPage(w, http.StatusInternalServerError, msg.H7003+err.Error(), r)

		return
	}
}

func RenderErrorPage(w http.ResponseWriter, statusCode int, msgCode string, r *http.Request) {
	SetSecureHTMLHeaders(w)

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

	w.WriteHeader(statusCode)
	RenderPage(w, tmplErrorHTML, data, r)
}

func RenderSinglePage(w http.ResponseWriter, tmplFile string, data interface{}, r *http.Request) {
	tmpl, err := template.ParseFiles(tmplFile)
	if err != nil {
		http.Error(w, "Template error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Template execution error: "+err.Error(), http.StatusInternalServerError)
		return
	}
}
