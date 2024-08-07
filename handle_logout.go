package uptime

import "net/http"

func handleLogout() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, cancelCookie())
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
	})
}
