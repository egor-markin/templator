package admin

import (
	"github.com/ReneKroon/ttlcache/v2"
	uuid "github.com/nu7hatch/gouuid"
	"html/template"
	"log"
	"net/http"
	"strings"
	"templator/handlers"
	"templator/services"
	"time"
)

// Based on
// https://www.sohamkamani.com/golang/session-based-authentication/
// https://github.com/sohamkamani/go-session-auth-example

var sessionTokenTTL time.Duration // Session token time to live

const sessionTokenCookieName = "session_token"

// Cache for session tokens
var sessionTokensCache ttlcache.SimpleCache

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		loginGetHandler(w)
	} else if r.Method == http.MethodPost {
		loginPostHandler(w, r)
	} else {
		http.Error(w, "This http method is not allowed", http.StatusMethodNotAllowed)
	}
}

func loginGetHandler(w http.ResponseWriter) {
	// GET: returning 'login.html'
	err := template.Must(template.ParseFiles("web/admin/login.html")).Execute(w, handlers.TemplatesPageData{ })
	if err != nil {
		log.Printf("Error occurred while trying to render HTML form: %s\n", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func loginPostHandler(w http.ResponseWriter, r *http.Request) {
	// Parsing form params
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	email := strings.ToLower(strings.TrimSpace(r.PostFormValue("email")))
	password := strings.TrimSpace(r.PostFormValue("password"))
	remember := strings.ToLower(strings.TrimSpace(r.PostFormValue("remember")))

	// Validating provided credentials
	if email != "admin@compasslabs.ru" || password != services.GetAdminPassword() {
		http.Error(w, "User name or password is unknown", http.StatusUnauthorized)
		return
	}

	// If user checked 'Remember me' checkbox, then saving its session as a cookie
	if remember == "on" {
		// Create a new random session token
		sessionToken, err := uuid.NewV4()
		if err != nil {
			log.Printf("Error occurred while trying to generate new UUID: %s\n", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		// Saving session token into the cache
		err = sessionTokensCache.Set(sessionToken.String(), nil)
		if err != nil {
			log.Printf("Error occurred while trying to save session token into the cache: %s\n", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		// Finally, we set the client cookie for "session_token" as the session token we just generated
		// we also set an expiry time of 120 seconds, the same as the sessionTokensCache
		http.SetCookie(w, &http.Cookie{
			Name:    sessionTokenCookieName,
			Value:   sessionToken.String(),
			Expires: time.Now().Add(sessionTokenTTL),
		})
	}


}

func init() {
	sessionTokenTTL = time.Duration(services.GetAdminSessionTokenTTL()) * time.Second

	sessionTokensCache = ttlcache.NewCache()
	err := sessionTokensCache.SetTTL(sessionTokenTTL)
	if err != nil {
		log.Printf("Unable to set TTL value for ttlcache: %s\n", err)
	}
}
