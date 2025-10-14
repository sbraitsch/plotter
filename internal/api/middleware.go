package api

import "net/http"

func AdminAuth(admins []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			key := r.Header.Get("X-Token")
			isAdmin := false
			for _, uuid := range admins {
				if key == uuid {
					isAdmin = true
				}
			}
			if !isAdmin {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
