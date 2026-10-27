package core_http_utils

import (
	"net/http"
)

func SetSessionCookie(w http.ResponseWriter, sessionID string, ttl int, isSecure bool) {
	http.SetCookie(w, &http.Cookie{
		Name:     "sessionID",
		Value:    sessionID,
		Path:     "/",
		MaxAge:   ttl,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   isSecure,
	})

}

func DeleteSessionCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:   "sessionID",
		MaxAge: -1,
		Path:   "/",
	})
}
