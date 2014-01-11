package middleware

import (
	"net/http"

	ctx "github.com/gorilla/context"
	"github.com/jordan-wright/gophish/auth"
)

// GetContext wraps each request in a function which fills in the context for a given request.
// This includes setting the User and Session keys and values as necessary for use in later functions.
func GetContext(handler http.Handler) http.Handler {
	// Set the context here
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set the context appropriately here.
		// Set the session
		session, _ := auth.Store.Get(r, "gophish")
		// Put the session in the context so that
		ctx.Set(r, "session", session)
		if id, ok := session.Values["id"]; ok {
			u, err := auth.GetUser(id.(int))
			if err != nil {
				ctx.Set(r, "user", nil)
			}
			ctx.Set(r, "user", u)
		} else {
			ctx.Set(r, "user", nil)
		}
		handler.ServeHTTP(w, r)
		// Save the session
		session.Save(r, w)
		// Remove context contents
		ctx.Clear(r)
	})
}

// RequireLogin is a simple middleware which checks to see if the user is currently logged in.
// If not, the function returns a 302 redirect to the login page.
func RequireLogin(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if u := ctx.Get(r, "user"); u != nil {
			handler.ServeHTTP(w, r)
		} else {
			http.Redirect(w, r, "/login", 302)
		}
	})
}
