package routers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	logger "testapp/internal/logs"
	"testapp/internal/storage/providers"
)

type router struct {
	pg *providers.Postgres_provider
}

func New(pg *providers.Postgres_provider) router {
	return router{pg}
}

// getbalance godoc
// @Summary Вывод баланса кошелька
// @Description Берет id кошелька и возвращает по нему баланс
// @Tags General
// @Accept json
// @Produce json
// @Param WALLET_UUID path int true "Wallet ID"
// @Success 200 {integer} models.Wallet.Amount "баланс кошелька"
// @Success 400 {object} string "Кошелек не найден"
// @Router /api/v1/wallet/{WALLET_UUID} [get]
func (router router) GetBalance(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("WALLET_UUID")

	balance, err := router.pg.GetBalance(id)
	if err != nil {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(balance)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(balance)
}

// getbalance godoc
// @Summary пополнение кошелька
// @Description по входным данным выдает операцию
// @Tags General
// @Accept json
// @Produce json
// @Param request body routers.UpdateBalance.Request true "query params"
// @Success 200 {object} string "операция успешна"
// @Success 400 {object} string "некорректный json/некорректное значение/недостаточно средств"
// @Success 404 {object} string "пользователь не найден"
// @Router /api/v1/wallet/update_balance [post]
func (router router) UpdateBalance(w http.ResponseWriter, r *http.Request) {
	type Request struct {
		ValletId      uint64 `json:"vallet_id"`
		OperationType string `json:"operation_type"`
		Amount        uint64 `json:"amount"`
	}

	var values Request

	err := json.NewDecoder(r.Body).Decode(&values)
	if err != nil {
		logger.Log.Info("некорректный json", slog.Any("error", err))
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("некорректный json")
		return
	}

	if values.Amount <= 0 || values.OperationType == "" || values.ValletId <= 0 {
		logger.Log.Info("функция UpdateBalance столкнулась с некорректным значением")
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("пришло некорректное значение")
		return
	}

	responce, err := router.pg.UpdateBalance(values.ValletId, values.OperationType, values.Amount)
	if err != nil {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(responce)
		return
	} else {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(responce)
		return
	}
}

// wallets godoc
// @Summary кошельки
// @Description выводит все кошельки
// @Tags General
// @Accept json
// @Produce json
// @Success 200 {object} string "операция успешна"
// @Success 404 {object} string "кршельков нет"
// @Router /api/v1/wallets [get]
func (router router) Wallets(w http.ResponseWriter, r *http.Request) {

	responce, err := router.pg.Wallets()
	if err != nil {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode("кошельков нет")
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(responce)
}
