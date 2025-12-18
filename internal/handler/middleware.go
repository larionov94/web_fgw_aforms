package handler

import (
	"errors"
	"fgw_web_aforms/internal/config"
	"fgw_web_aforms/internal/handler/http_err"
	"fgw_web_aforms/internal/service"
	"fgw_web_aforms/pkg/common"
	"fgw_web_aforms/pkg/common/msg"
	"fmt"
	"html/template"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/sessions"
)

const (
	pathToDefault        = "/"
	prefixTmplPerformers = "web/html/"
	tmplForceLogoutHTML  = "force_logout.html"
	maxLifeSession       = 4 * time.Hour
	expiresCache         = 60 * time.Minute
)

type AuthMiddleware struct {
	store        *sessions.CookieStore
	sessName     string
	performerKey string
	roleKey      string
	logg         *common.Logger
	userCache    map[int]*UserSession
	cacheMu      sync.RWMutex
}

type UserSession struct {
	PerformerFIO string
	RoleName     string
	Expires      time.Time
}

func NewAuthMiddleware(store *sessions.CookieStore, logg *common.Logger) *AuthMiddleware {
	return &AuthMiddleware{
		store:        store,
		sessName:     config.GetSessionName(),
		performerKey: config.SessionPerformerKey,
		roleKey:      config.SessionRoleKey,
		logg:         logg,
	}
}

// GetUserData получаем данные пользователя.
func (m *AuthMiddleware) GetUserData(r *http.Request, performerService service.PerformerUseCase, roleService service.RoleUseCase) (string, int, string, error) {
	performerId, ok1 := m.GetPerformerId(r)
	roleId, ok2 := m.GetRoleId(r)
	if !ok1 || !ok2 {
		m.logg.LogE(msg.E3103, nil)

		return "", 0, "", errors.New("пользователь не авторизован")
	}

	var performerFIO string
	var roleName string

	// 1. Проверяем кеш.
	m.cacheMu.RLock()
	if cached, exists := m.userCache[performerId]; exists && time.Now().Before(cached.Expires) {
		performerFIO = cached.PerformerFIO
		roleName = cached.RoleName
		m.cacheMu.RUnlock()

		return performerFIO, performerId, roleName, nil
	}
	m.cacheMu.RUnlock()

	// 2. Загружаем данные.
	ctx := r.Context()

	performer, err := performerService.FindByIdPerformer(ctx, performerId)
	if err != nil {
		m.logg.LogE(msg.E3206, err)

		return "", performerId, "", err
	}

	role, err := roleService.FindRoleById(ctx, roleId)
	if err != nil {
		m.logg.LogE(msg.E3206, err)

		return performerFIO, performerId, "", err
	}

	// 3. Сохраняем в кеш.
	m.cacheMu.Lock()
	if m.userCache == nil {
		m.userCache = make(map[int]*UserSession)
	}
	m.userCache[performerId] = &UserSession{
		PerformerFIO: performer.FIO,
		RoleName:     role.Name,
		Expires:      time.Now().Add(expiresCache),
	}
	m.cacheMu.Unlock()

	return performer.FIO, performerId, role.Name, nil
}

// RequireAuth - основной middleware для проверки аутентификации.
func (m *AuthMiddleware) RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Устанавливаем заголовки безопасности для всех защищенных запросов.
		m.setSecurityHeaders(w)

		// Получаем сессию с помощью безопасного метода.
		session, err := m.getSecureSession(r)
		if err != nil {
			m.forceLogoutAndRedirect(w, r, msg.H7013)

			return
		}

		if session == nil {
			m.forceLogoutAndRedirect(w, r, msg.H7012)

			return
		}

		// Проверяем аутентификацию.
		if auth, ok := session.Values[config.SessionAuthPerformer].(bool); !ok || !auth {
			m.forceLogoutAndRedirect(w, r, msg.H7011)

			return
		}

		// Проверяем время жизни сессии.
		if m.isSessionExpired(session) {
			m.forceLogoutAndRedirect(w, r, msg.H7010)

			return
		}

		// Обновляем активность сессии.
		m.updateSessionActivity(session, w, r)

		// Для HTML-ответов добавляем скрипт управления историей.
		if r.Header.Get("Accept") == "text/html" {
			m.addHistoryManagementScript(w)
		}

		next.ServeHTTP(w, r)
	}
}

// RequireRole - middleware для проверки ролей.
func (m *AuthMiddleware) RequireRole(requireRoles []int, next http.HandlerFunc) http.HandlerFunc {
	allowedRoles := make(map[int]bool)
	for _, role := range requireRoles {
		allowedRoles[role] = true
	}

	return m.RequireAuth(func(w http.ResponseWriter, r *http.Request) {
		session, err := m.store.Get(r, m.sessName)
		if err != nil {
			http.Redirect(w, r, pathToDefault, http.StatusFound)
			return
		}

		performerRole, ok := session.Values[m.roleKey].(int)
		if !ok {
			m.forceLogoutAndRedirect(w, r, msg.H7014)
			return
		}

		if !allowedRoles[performerRole] {
			http_err.SendErrorHTTP(w, http.StatusForbidden, msg.H7015, m.logg, r)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// GetPerformerId - получение ID пользователя.
func (m *AuthMiddleware) GetPerformerId(r *http.Request) (int, bool) {
	session, err := m.store.Get(r, m.sessName)
	if err != nil {
		return 0, false
	}

	performerId, ok := session.Values[m.performerKey].(int)
	return performerId, ok
}

// GetRoleId - получение ID роли.
func (m *AuthMiddleware) GetRoleId(r *http.Request) (int, bool) {
	session, err := m.store.Get(r, m.sessName)
	if err != nil {
		return 0, false
	}

	performerRole, ok := session.Values[m.roleKey].(int)
	return performerRole, ok
}

// getSecureSession - безопасное получение сессии с валидацией.
func (m *AuthMiddleware) getSecureSession(r *http.Request) (*sessions.Session, error) {
	session, err := m.store.Get(r, m.sessName)
	if err != nil {
		return nil, err
	}

	// Дополнительная валидация куки.
	if session.IsNew {
		return nil, nil
	}

	return session, nil
}

// isSessionExpired - проверка истечения срока действия сессии.
func (m *AuthMiddleware) isSessionExpired(session *sessions.Session) bool {
	if createdAt, ok := session.Values["created_at"].(int64); ok {
		createTime := time.Unix(createdAt, 0)

		maxAge := maxLifeSession
		if customMaxAge, ok := session.Values["max_age"].(int); ok {
			maxAge = time.Duration(customMaxAge) * time.Second
		}

		return time.Since(createTime) > maxAge
	}

	return true
}

// updateSessionActivity - обновление времени активности.
func (m *AuthMiddleware) updateSessionActivity(session *sessions.Session, w http.ResponseWriter, r *http.Request) {
	now := time.Now()
	session.Values["last_activity"] = now.Unix()

	// Устанавливаем куку с коротким временем жизни для браузера.
	if cookie, err := r.Cookie("activity_check"); err != nil || cookie.Value != "active" {
		http.SetCookie(w, &http.Cookie{
			Name:     "activity_check",
			Value:    "active",
			Path:     pathToDefault,
			MaxAge:   -1,
			HttpOnly: false,
			Secure:   true,
			SameSite: http.SameSiteStrictMode,
		})
	}

	if err := session.Save(r, w); err != nil {
		return
	}
}

// setSecurityHeaders - установка заголовков безопасности.
func (m *AuthMiddleware) setSecurityHeaders(w http.ResponseWriter) {
	w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, private, max-age=0")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("X-Frame-Options", "DENY")
}

// addHistoryManagementScript - добавление скрипта управления историей.
func (m *AuthMiddleware) addHistoryManagementScript(w http.ResponseWriter) {
	w.Header().Add("X-History-Control", "no-cache")
}

// forceLogoutAndRedirect - принудительный выход и редирект с очисткой истории.
func (m *AuthMiddleware) forceLogoutAndRedirect(w http.ResponseWriter, r *http.Request, reason string) {
	m.logg.LogW(fmt.Sprintf("%s: %s", msg.H8000, reason))

	// Уничтожаем сессию.
	if session, err := m.store.Get(r, m.sessName); err == nil {

		// Очищаем сессию.
		session.Options.MaxAge = -1
		for key := range session.Values {
			delete(session.Values, key)
		}
		if err = session.Save(r, w); err != nil {
			return
		}
	}

	// Устанавливаем заголовки no-cache.
	m.setSecurityHeaders(w)

	// Загружаем HTML шаблон
	tmpl, err := template.ParseFiles(prefixTmplPerformers + tmplForceLogoutHTML)
	if err != nil {
		http.Error(w, msg.H7016, http.StatusInternalServerError)

		return
	}

	// Устанавливаем заголовок ответа.
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	data := struct {
		Reason    string
		Timestamp string
	}{
		Reason:    reason,
		Timestamp: time.Now().Format("02.01.2006 15:04:05"),
	}

	if err = tmpl.Execute(w, data); err != nil {
		http.Error(w, msg.H7017, http.StatusInternalServerError)

		return
	}
}
