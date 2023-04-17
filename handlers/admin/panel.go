package admin

import (
	"fmt"
	"github.com/ReneKroon/ttlcache/v2"
	"net/http"
)

func PanelHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		panelGetHandler(w, r)
	} else {
		http.Error(w, "This http method is not allowed", http.StatusMethodNotAllowed)
	}
}

func panelGetHandler(w http.ResponseWriter, r *http.Request) {
	// We can obtain the session token from the requests cookies, which come with every request
	c, err := r.Cookie(sessionTokenCookieName)
	if err != nil {
		if err == http.ErrNoCookie {
			//http.RedirectHandler()
			// If the cookie is not set, return an unauthorized status
			http.Error(w, "Please sign in before continuing", http.StatusUnauthorized)
			return
		}
		// For any other type of error, return a bad request status
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	sessionToken := c.Value

	// We then get the name of the user from our cache, where we set the session token
	_, err = sessionTokensCache.Get(sessionToken)
	if err != nil {
		if err == ttlcache.ErrNotFound {
			// If the session token is not present in cache, return an unauthorized error
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		// If there is an error fetching from cache, return an internal server error status
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Finally, return the welcome message to the user
	w.Write([]byte(fmt.Sprintf("Welcome!")))
}