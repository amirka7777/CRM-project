package main

import (
	"log"
	"net/http"

	"github.com/yukay/CRM/internal/database"
	"github.com/yukay/CRM/internal/handlers"
	"github.com/yukay/CRM/internal/middleware"
	"github.com/yukay/CRM/internal/repository"
	"github.com/yukay/CRM/internal/service"
)

func main() {

	secretKey := []byte("super_secret_key")
	db, err := database.InitDB("base.db")
	if err != nil {
		log.Fatalf("произошла критическая ошибка при инициализации БД: %v", err)
	}

	defer func() {
		if err := db.Close(); err != nil {
			log.Println("Ошибка при закрытии бд: ", err)
		} else {
			log.Println("База данных успешно закрылась")
		}
	}()

	repo := repository.NewRepository(db)
	err = database.CreateUsersTable(db)
	if err != nil {
		log.Fatalf("произошла критическая ошибка при создании таблицы: %v", err)
	}
	err = database.CreateDealsTable(db)
	if err != nil {
		log.Fatalf("произошла критическая ошибка при создании таблицы: %v", err)
	}
	err = database.CreateContactsTable(db)
	if err != nil {
		log.Fatalf("произошла критическая ошибка при создании таблицы: %v", err)
	}

	service := service.NewService(repo, secretKey)
	handler := handlers.NewHandlerSupport(service)

	http.HandleFunc("/auth/register", handler.CreateUser)
	http.HandleFunc("/auth/login", handler.LoginUser)
	http.Handle("/me", middleware.AuthMiddleware(secretKey)(http.HandlerFunc(handler.Me)))
	http.Handle("/contacts", middleware.AuthMiddleware(secretKey)(http.HandlerFunc(handler.CreateContact)))

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Println("ошибка при запуске сервера")
	}
	log.Println("Сервер успешно запущен")
}
