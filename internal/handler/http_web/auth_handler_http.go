package http_web

import (
	"fgw_web_aforms/internal/config"
	"fgw_web_aforms/internal/handler"
	"fgw_web_aforms/internal/handler/http_err"
	"fgw_web_aforms/internal/handler/page"
	"fgw_web_aforms/internal/service"
	"fgw_web_aforms/pkg/common"
	"fgw_web_aforms/pkg/common/msg"
	"fgw_web_aforms/pkg/convert"
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/sessions"
)

const (
	tmplRedirectHTML      = "redirect.html"
	tmplAuthHTML          = "auth.html"
	tmplProductionHTML    = "productions.html"
	tmplProductionAddHTML = "production_add.html"

	urlAForms             = "/aforms"
	urlAuth               = "/auth"
	urlLogin              = "/login"
	urlLogoutTempRedirect = "/logout-temp-redirect"
	urlTempRedirect       = "/temp-redirect"
	pathToDefault         = "/"
	tmplStartPageHTML     = "index.html"
)

const (
	RedirectDelayFast    = 100  // 0.1 секунда
	RedirectDelayNormal  = 300  // 0.3 секунды
	FallbackDelayDefault = 3000 // 3 секунды
)

type AuthHandlerHTML struct {
	performerService service.PerformerUseCase
	roleService      service.RoleUseCase
	logg             *common.Logger
	authMiddleware   *handler.AuthMiddleware
}

type RedirectData struct {
	Title           string
	Message         string
	NoScriptMessage string
	TargetURL       string
	CurrentURL      string
	TempURL         string
	Delay           int
	FallbackDelay   int
	ClearHistory    bool
	AddTempState    bool
}

func NewAuthHandlerHTML(
	performerService service.PerformerUseCase,
	roleService service.RoleUseCase,
	logg *common.Logger,
	authMiddleware *handler.AuthMiddleware) *AuthHandlerHTML {

	return &AuthHandlerHTML{
		performerService: performerService,
		roleService:      roleService,
		logg:             logg,
		authMiddleware:   authMiddleware,
	}
}

func (a *AuthHandlerHTML) ServerHTTPRouter(mux *http.ServeMux) {
	mux.HandleFunc("/", a.ShowAuthForm)
	mux.HandleFunc("/login", a.LoginPage)
	mux.HandleFunc("/auth", a.AuthPerformerHTML)
	mux.HandleFunc("/logout", a.Logout)
	mux.HandleFunc("/aforms", a.authMiddleware.RequireAuth(a.authMiddleware.RequireRole([]int{0, 4, 5}, a.StartPage)))
}

func (a *AuthHandlerHTML) StartPage(w http.ResponseWriter, r *http.Request) {
	performerData, err := a.authMiddleware.GetUserData(r, a.performerService, a.roleService)
	if err != nil {
		a.redirectToLoginWithHistoryClear(w, r)

		return
	}

	data := page.NewDataPage(
		"Панель форма комплектов", "dashboard", performerData,
		nil,
		nil,
		nil,
		false)

	page.RenderPages(w, tmplStartPageHTML, data, r, tmplProductionHTML, tmplProductionAddHTML)
}

func (a *AuthHandlerHTML) ShowAuthForm(w http.ResponseWriter, r *http.Request) {
	session, err := config.Store.Get(r, config.GetSessionName())
	if err == nil {
		if auth, ok := session.Values[config.SessionAuthPerformer].(bool); ok && auth {
			a.safeRedirectBasedOnRole(w, r, session)
			return
		}
	}

	a.LoginPage(w, r)
}

func (a *AuthHandlerHTML) LoginPage(w http.ResponseWriter, r *http.Request) {
	page.SetSecureHTMLHeaders(w)

	if r.Method != http.MethodGet {
		http_err.SendErrorHTTP(w, http.StatusMethodNotAllowed, "", a.logg, r)

		return
	}

	errorMsg := r.URL.Query().Get("error")

	data := struct {
		ErrorMessage string
	}{
		ErrorMessage: errorMsg,
	}

	page.RenderPage(w, tmplAuthHTML, data, r)
}

func (a *AuthHandlerHTML) Logout(w http.ResponseWriter, r *http.Request) {
	session, err := config.Store.Get(r, config.GetSessionName())
	if err != nil {
		a.sendLogoutPageWithHistoryClear(w, r)
		return
	}

	if token, ok := session.Values["session_token"].(string); ok {
		if mw, ok := interface{}(a.authMiddleware).(interface{ RemoveSessionToken(token string) }); ok {
			mw.RemoveSessionToken(token)
		}
	}

	for key := range session.Values {
		delete(session.Values, key)
	}

	session.Options.MaxAge = -1
	session.Options.HttpOnly = true
	session.Options.Secure = true
	session.Options.SameSite = http.SameSiteStrictMode

	if err = session.Save(r, w); err != nil {
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     config.GetSessionName(),
		Value:    "",
		Path:     pathToDefault,
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})

	a.sendLogoutPageWithHistoryClear(w, r)
}

func (a *AuthHandlerHTML) AuthPerformerHTML(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http_err.SendErrorHTTP(w, http.StatusMethodNotAllowed, "", a.logg, r)
		return
	}

	if err := r.ParseForm(); err != nil {
		page.RenderErrorPage(w, http.StatusBadRequest, msg.H7007, r)
		return
	}

	performerIdStr := r.FormValue("performerId")
	performerPass := r.FormValue("performerPassword")

	if performerIdStr == "" || performerPass == "" {
		page.RenderErrorPage(w, http.StatusUnauthorized, msg.E3211, r)
		return
	}

	performerId := convert.ConvStrToInt(performerIdStr)

	authResult, err := a.performerService.AuthPerformer(r.Context(), performerId, performerPass)
	if err != nil {
		if authResult != nil && !authResult.Success {
			http.Redirect(w, r, "/login?error="+url.QueryEscape(authResult.Message), http.StatusFound)
		} else {
			http.Redirect(w, r, "/login?error="+url.QueryEscape(msg.H7005), http.StatusFound)
		}
		return
	}

	if authResult.Success {
		err := a.createSecureSession(w, r, performerId, authResult.Performer.IdRoleAForms)
		if err != nil {
			page.RenderErrorPage(w, http.StatusInternalServerError, "Ошибка создания сессии", r)
			return
		}

		a.sendLoginSuccessPage(w, r, authResult.Performer.IdRoleAForms)
	} else {
		http.Redirect(w, r, "/login?error="+url.QueryEscape(authResult.Message), http.StatusFound)
	}
}

// НОВЫЙ МЕТОД: safeRedirectBasedOnRole с использованием общего шаблона
func (a *AuthHandlerHTML) safeRedirectBasedOnRole(w http.ResponseWriter, r *http.Request, session *sessions.Session) {
	target := urlAForms

	data := RedirectData{
		Title:           "Перенаправление",
		Message:         "Вы уже авторизованы. Выполняется безопасное перенаправление...",
		NoScriptMessage: "Включите JavaScript для безопасного перехода.",
		TargetURL:       target,
		CurrentURL:      r.URL.Path,
		TempURL:         urlTempRedirect,
		Delay:           RedirectDelayFast,
		FallbackDelay:   FallbackDelayDefault,
		ClearHistory:    true,
		AddTempState:    false,
	}

	a.renderRedirectPage(w, r, data)
}

func (a *AuthHandlerHTML) renderRedirectPage(w http.ResponseWriter, r *http.Request, data RedirectData) {
	if data.Title == "" {
		data.Title = "Перенаправление"
	}
	if data.Message == "" {
		data.Message = "Выполняется безопасное перенаправление..."
	}
	if data.NoScriptMessage == "" {
		data.NoScriptMessage = "Включите JavaScript для безопасного перехода."
	}
	if data.CurrentURL == "" {
		data.CurrentURL = r.URL.Path
	}
	if data.Delay == 0 {
		data.Delay = RedirectDelayNormal
	}
	if data.FallbackDelay == 0 {
		data.FallbackDelay = FallbackDelayDefault
	}

	page.SetSecureHTMLHeaders(w)
	page.RenderPage(w, tmplRedirectHTML, data, r)
}

// Обновленный sendLoginSuccessPage
func (a *AuthHandlerHTML) sendLoginSuccessPage(w http.ResponseWriter, r *http.Request, roleId int) {
	target := urlAForms

	data := RedirectData{
		Title:           "Успешный вход",
		Message:         "Вход выполнен успешно. Выполняется безопасное перенаправление...",
		NoScriptMessage: "Включите JavaScript для безопасного перехода.",
		TargetURL:       target,
		CurrentURL:      urlAuth,
		TempURL:         urlLogoutTempRedirect,
		Delay:           RedirectDelayNormal,
		FallbackDelay:   2000,
		ClearHistory:    true,
		AddTempState:    true,
	}

	a.renderRedirectPage(w, r, data)
}

// Обновленный sendLogoutPageWithHistoryClear
func (a *AuthHandlerHTML) sendLogoutPageWithHistoryClear(w http.ResponseWriter, r *http.Request) {
	data := RedirectData{
		Title:           "Выход из системы",
		Message:         "Вы успешно вышли из системы. Выполняется безопасное перенаправление на страницу входа...",
		NoScriptMessage: "Включите JavaScript для безопасного выхода.",
		TargetURL:       urlLogin,
		CurrentURL:      r.URL.Path,
		TempURL:         urlLogoutTempRedirect,
		Delay:           RedirectDelayNormal,
		FallbackDelay:   FallbackDelayDefault,
		ClearHistory:    true,
		AddTempState:    true,
	}

	a.renderRedirectPage(w, r, data)
}

func (a *AuthHandlerHTML) redirectToLoginWithHistoryClear(w http.ResponseWriter, r *http.Request) {
	a.sendLogoutPageWithHistoryClear(w, r)
}

func (a *AuthHandlerHTML) createSecureSession(w http.ResponseWriter, r *http.Request, performerId, roleId int) error {
	session, _ := config.Store.Get(r, config.GetSessionName())

	token := config.GenerateSessionToken()

	session.Values[config.SessionAuthPerformer] = true
	session.Values[config.SessionPerformerKey] = performerId
	session.Values[config.SessionRoleKey] = roleId
	session.Values["session_token"] = token
	session.Values["created_at"] = time.Now().Unix()
	session.Values["last_activity"] = time.Now().Unix()

	session.Options = &sessions.Options{
		Path:     pathToDefault,
		MaxAge:   1800,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	}

	page.SetSecureHTMLHeaders(w)

	return session.Save(r, w)
}
