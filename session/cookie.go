//package session handles the cookie & session management
package session

import (
	"golive/envvar"
	"golive/logger"
	"net/http"
	"time"

	uuid "github.com/satori/go.uuid"
)

//CheckCookie returns cookie status
func CheckCookie(r *http.Request) bool {
	_, err := r.Cookie(envvar.CookieName())
	if err != nil {
		return false
	}
	return true
}

//SetCookie creates new cookie for the client
func SetCookie(w http.ResponseWriter, r *http.Request) {
	cookieName := envvar.CookieName()
	cookie, err := r.Cookie(cookieName)
	if err != nil {
		cookie := &http.Cookie{
			Name:     cookieName,
			Value:    (uuid.NewV4()).String(),
			HttpOnly: true,
			Path:     "/",
			Domain:   envvar.HostAddress(),
			Secure:   true,
		}
		http.SetCookie(w, cookie)
		logger.Info.Println("set cookie : " + cookie.Value + "-" + cookieName)
		return
	}
	_, found := Get(r)
	if found {
		Refresh(r)
		logger.Info.Println("session refresh: " + cookie.Value)
		return
	}
	logger.Info.Println(cookie.Value + " already set")

	return
}

//DeleteCookie removes cookie from the client
//by changing the expires & maxAge value
func DeleteCookie(w http.ResponseWriter, r *http.Request) {
	cookieName := envvar.CookieName()
	cookie, err := r.Cookie(cookieName)
	if err == nil {
		logger.Info.Println("remove cookie : " + cookie.Value)
		cookie = &http.Cookie{
			Name:     cookieName,
			Expires:  time.Now().AddDate(-1, 0, 0),
			Value:    "",
			MaxAge:   -1,
			Path:     "/",
			HttpOnly: true,
			Secure:   true,
		}
		http.SetCookie(w, cookie)
	}
	return
}
