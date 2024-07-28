package handlers

import (
	"context"
	"net/http"
	"text/template"
	"time"

	"balance-tracker/models"
	"balance-tracker/services"
	"balance-tracker/utils"
)

type AuthHandler struct {
	authService services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{*authService}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	if username == "" || password == "" {
		h.renderTemplate(w, "login", map[string]string{"Error": "Username and password are required"})
		return
	}

	token, err := h.authService.Login(username, password)
	if err != nil {
		h.renderTemplate(w, "login", map[string]string{"Error": err.Error()})
		return
	}

	// Set the token as a cookie
	cookie := &http.Cookie{
		Name:  "token",
		Value: token,
		Path:  "/",
	}

	http.SetCookie(w, cookie)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *AuthHandler) renderTemplate(w http.ResponseWriter, name string, data interface{}) {
	tmpl, err := template.ParseFiles("templates/" + name + ".html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	user := models.User{
		Username: r.Form.Get("username"),
		Password: r.Form.Get("password"),
	}

	err := h.authService.Register(user)
	if err != nil {
		data := struct {
			Error string
		}{
			Error: err.Error(),
		}
		tmpl, err := template.ParseFiles("templates/register.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, data)
		return
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusBadRequest)
		return
	}

	cookie, err := r.Cookie("token")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	// Get the token value from the cookie
	tokenString := cookie.Value
	err = h.authService.Logout(tokenString)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Clear the cookie
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   "",
		Expires: time.Unix(0, 0),
	})

	http.Redirect(w, r, "/login", http.StatusFound)
	return
}

func (h *AuthHandler) AuthMiddleware() func(next http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			// Get the cookie by name
			cookie, err := r.Cookie("token")
			if err != nil {
				http.Redirect(w, r, "/login", http.StatusFound)
				return
			}

			// Get the token value from the cookie
			tokenString := cookie.Value

			// Check if the token exists in the database
			if !h.authService.TokenValid(tokenString) {
				// Delete the token cookie
				http.SetCookie(w, &http.Cookie{
					Name:    "token",
					Value:   "",
					Expires: time.Unix(0, 0),
				})
				http.Redirect(w, r, "/login", http.StatusFound)
				return
			}

			// Parse the JWT token
			claims, err := utils.ParseToken(tokenString)
			if err != nil {
				http.Redirect(w, r, "/login", http.StatusFound)
				return
			}

			// Store the UserID in the request context
			ctx := context.WithValue(r.Context(), "userID", claims.UserID)

			// Call the next handler function with the updated context
			next(w, r.WithContext(ctx))
		}
	}
}
