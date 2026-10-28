package core_http_utils

import (
	"net/http"
	"time"
)

func SetSessionCookie(w http.ResponseWriter, sessionID string, ttl int, isSecure bool) {
	ttl = 86400
	expires := time.Now().Add(24 * time.Hour)
	http.SetCookie(w, &http.Cookie{
		Name:     "sessionID",
		Value:    sessionID,
		Path:     "/",
		MaxAge:   ttl,
		Expires:  expires.UTC(),
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
