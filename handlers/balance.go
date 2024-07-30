package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"path"
	"strconv"
	"text/template"

	"balance-tracker/models"
	"balance-tracker/services"

	"github.com/go-chi/chi/v5"
)

type BalanceHandler struct {
	balanceService services.BalanceService
}

func NewBalanceHandler(balanceService *services.BalanceService) *BalanceHandler {
	return &BalanceHandler{*balanceService}
}

func (h *BalanceHandler) GetBalances(w http.ResponseWriter, r *http.Request) {
	balances, err := h.balanceService.GetBalances()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(balances)
}

func (h *BalanceHandler) GetBalance(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	balance, err := h.balanceService.GetBalance(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(balance)
}

func (h *BalanceHandler) UpdateBalance(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var balance models.Balance
	err := json.NewDecoder(r.Body).Decode(&balance)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.balanceService.UpdateBalance(id, balance)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *BalanceHandler) CreateBalance(w http.ResponseWriter, r *http.Request) {
	// Get the UserID from the claims
	userID := r.Context().Value("userID").(int)

	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	amount, err := strconv.ParseFloat(r.Form.Get("amount"), 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	err = h.balanceService.CreateBalance(models.Balance{
		UserID: userID,
		Amount: amount,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	balance, err := h.balanceService.GetLastBalanceByUserID(userID)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl := template.Must(template.ParseFiles("templates/components/balanceCard.html"))
	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, balance)

}

func (h *BalanceHandler) DeleteBalance(w http.ResponseWriter, r *http.Request) {
	id := path.Base(r.URL.Path)
	if id == "" {
		log.Println("Missing id parameter")
		http.Error(w, "Missing id parameter", http.StatusBadRequest)
		return
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid id", http.StatusBadRequest)
		return
	}

	err = h.balanceService.DeleteBalance(idInt)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *BalanceHandler) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int)

	earn, err := strconv.Atoi(r.FormValue("earn"))
	if err != nil {
		http.Error(w, "Invalid earn value", http.StatusBadRequest)
		return
	}

	expense, err := strconv.Atoi(r.FormValue("expense"))
	if err != nil {
		http.Error(w, "Invalid expense value", http.StatusBadRequest)
		return
	}

	amount := earn - expense

	balance, err := h.balanceService.CreateNewTransactionByID(userID, amount)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the balanceCard.html template
	tmpl, err := template.ParseFiles("templates/components/balanceCard.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := struct {
		Amount    float64
		CreatedAt string
		ID        int
	}{
		Amount:    balance.Amount,
		CreatedAt: balance.CreatedAt,
		ID:        balance.ID,
	}
	log.Println(data)

	// if err := tmpl.Execute(w, data); err != nil {
	// 	log.Println(err)
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }
	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, data)
}
