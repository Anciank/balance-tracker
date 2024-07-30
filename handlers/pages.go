// handlers/page.go

package handlers

import (
	"balance-tracker/models"
	"balance-tracker/services"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type PageHandler struct {
	template       *template.Template
	balanceService *services.BalanceService
}

func NewPageHandler(balanceService *services.BalanceService) *PageHandler {
	tmpl, err := template.ParseFiles("templates/index.html", "templates/login.html", "templates/register.html")
	if err != nil {
		log.Fatal(err)
		return nil
	}
	return &PageHandler{
		template:       tmpl,
		balanceService: balanceService,
	}
}

func (h *PageHandler) HandleLoginPage(w http.ResponseWriter, r *http.Request) {
	err := h.template.ExecuteTemplate(w, "login.html", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *PageHandler) HandleRegisterPage(w http.ResponseWriter, r *http.Request) {
	err := h.template.ExecuteTemplate(w, "register.html", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *PageHandler) HandleIndexPage(w http.ResponseWriter, r *http.Request) {
	type BalancePage struct {
		Balances []models.Balance
	}

	// Retrieve the UserID from the request context
	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		http.Error(w, "Invalid user ID", http.StatusInternalServerError)
		return
	}

	balances, err := h.balanceService.GetBalancesByUserID(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	balancePage := BalancePage{
		Balances: balances,
	}

	err = h.template.ExecuteTemplate(w, "index.html", balancePage)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *PageHandler) HandleHtmxServe(w http.ResponseWriter, r *http.Request) {
	wd, err := os.Getwd()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	filePath := filepath.Join(wd, "public", "htmx.min.js")
	http.ServeFile(w, r, filePath)
}

func (h *PageHandler) HandleTailwindServe(w http.ResponseWriter, r *http.Request) {
	wd, err := os.Getwd()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	filePath := filepath.Join(wd, "public", "tailwind.js")
	http.ServeFile(w, r, filePath)
}

func (h *PageHandler) HandleStaticServe(w http.ResponseWriter, r *http.Request) {
	// Get the file path from the URL
	filePath := r.URL.Path

	// Remove the /static/ prefix from the file path
	filePath = strings.TrimPrefix(filePath, "/static/")

	// Define the two directories to search for the file
	dirs := []string{"templates/components/", "public/"}

	// Loop through the directories and try to find the file
	for _, dir := range dirs {
		file, err := os.Open(dir + filePath)
		if err == nil {
			// File found, serve it
			defer file.Close()
			http.ServeFile(w, r, file.Name())
			return
		}
	}

	// File not found in either directory, return 404
	http.Error(w, "File not found", http.StatusNotFound)
}
