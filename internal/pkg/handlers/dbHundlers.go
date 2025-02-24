package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"testapp/internal/pkg/models"
)

func (h handler) Hello(w http.ResponseWriter, r *http.Request) {
	results := "hello"
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&results)
}

func (h handler) GetBalance(w http.ResponseWriter, r *http.Request) {

	id := r.PathValue("WALLET_UUID")
	var wallet = models.Wallet{}

	result := h.DB.Find(&wallet, id)
	if result.Error != nil {
		log.Println(result.Error)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadGateway)
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&wallet.Amount)
}

func (h handler) UpdateBalance(w http.ResponseWriter, r *http.Request) {

	type Jsonvalues struct {
		ValletId      uint64 `json:"vallet_id"`
		OperationType string `json:"operation_type"`
		Amount        uint64 `json:"amount"`
	}

	var values Jsonvalues
	err := json.NewDecoder(r.Body).Decode(&values)
	if err != nil {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("некорректный json")
		return
	}

	log.Println(values)
	if values.Amount == 0 || values.OperationType == "" || values.ValletId == 0 {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("пришло пустое значение")
		return
	}
	wallet := models.Wallet{}

	result := h.DB.First(&wallet, values.ValletId)
	if result.Error != nil {
		log.Println(result.Error)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("пользователь не найден")
		return
	}

	switch {
	case values.OperationType == "DEPOSIT":
		wallet.Amount = wallet.Amount + values.Amount
		h.DB.Save(&wallet)
		log.Println("счет пополнен")
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode("счет пополнен")
		return
	case values.OperationType == "WITHDRAW":
		if values.Amount > wallet.Amount {

			log.Println("недостаточно средств")
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode("недостаточно средств")
			return
		} else {
			wallet.Amount = wallet.Amount - values.Amount
			h.DB.Save(&wallet)
			log.Println("операция успешна")
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode("операция успешна")
			return
		}
	}

}
