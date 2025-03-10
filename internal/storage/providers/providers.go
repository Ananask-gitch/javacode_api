package providers

import (
	"errors"
	"log/slog"
	logger "testapp/internal/logs"
	"testapp/internal/models"

	"gorm.io/gorm"
)

type Postgres_provider struct {
	db *gorm.DB
}

func New(db *gorm.DB) Postgres_provider {
	return Postgres_provider{db}
}

func (pg Postgres_provider) GetBalance(id string) (any, error) {
	logger.Log.Info("начало получения баланса кошелька: [*id = UUID кошелбка*]", slog.String("id:", id))
	slog.Info("начало получения баланса кошелька: [*id = UUID кошелбка*]", slog.String("id:", id))
	var wallet = models.Wallet{}
	result := pg.db.Take(&wallet, id)
	if result.Error != nil {
		responce := "кошелек не найден"
		logger.Log.Info(responce)
		slog.Info(responce)
		return responce, result.Error
	}

	return wallet.Amount, nil
}

func (pg Postgres_provider) UpdateBalance(valletId uint64, operationType string, Amount uint64) (any, error) {

	wallet := models.Wallet{}
	result := pg.db.Take(&wallet, valletId)
	if result.Error != nil {
		logger.Log.Info("кошелек не найден")
		responce := "Кошелек не найден"
		return responce, result.Error
	}

	switch {
	case operationType == "DEPOSIT":
		wallet.Amount = wallet.Amount + Amount
		pg.db.Save(&wallet)
		logger.Log.Info("функция UpdateBalance успешно выполнена")
		return nil, nil
	case operationType == "WITHDRAW":
		if Amount > wallet.Amount {
			responce := "недостаточно средств"
			logger.Log.Info("недостаточно средств")
			return responce, errors.New("недостаточно средств")
		} else {
			wallet.Amount = wallet.Amount - Amount
			pg.db.Save(&wallet)
			logger.Log.Info("функция UpdateBalance успешно выполнена")
			return nil, nil
		}
	}
	return nil, errors.New("что-то пошло не так")
}

func (pg Postgres_provider) Wallets() (any, error) {
	var wallet []models.Wallet
	result := pg.db.Find(&wallet)
	if result.Error != nil {
		return nil, errors.New("не получилось")
	}
	return &wallet, nil
}
