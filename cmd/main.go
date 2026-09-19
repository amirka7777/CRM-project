package main

import (
	"log"
	"net/http"

	"github.com/yukay/CRM/internal/handlers"
)

func main() {

	http.HandleFunc("/", handlers.Ping)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Println("ошибка при запуске сервера")
	}
	log.Println("Сервер успешно запущен")
}
