package storage

import (
	"fmt"
	"log"
	app "testapp/internal/configs"
	"testapp/internal/pkg/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func DatabaseInit(config app.Config) *gorm.DB {

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.DBHost, config.DBPort, config.DBUser, config.DBPassword, config.DBName)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Printf("Ошибка подключения базы данных: %v", err)
		dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s sslmode=disable",
			config.DBHost, config.DBPort, config.DBUser, config.DBPassword)
		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Fatalf("Ошибка подключения к постгрес %v", err)
		}
		createDatabaseCommand := fmt.Sprintf("CREATE DATABASE %s", config.DBName)
		db.Exec(createDatabaseCommand)
	}

	err = db.AutoMigrate(&models.Wallet{})
	if err != nil {
		log.Fatalf("Ошибка миграции: %v", err)
	}

	log.Println("База данных успешно создана и миграция завершена.")

	return db
}
