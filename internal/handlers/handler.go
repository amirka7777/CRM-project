package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/yukay/CRM/internal/middleware"
	"github.com/yukay/CRM/internal/models"
	"github.com/yukay/CRM/internal/service"
)

type HandlerSupport struct {
	Serv *service.ServiceSupport
}

func NewHandlerSupport(serv *service.ServiceSupport) *HandlerSupport {
	return &HandlerSupport{Serv: serv}
}

func (h *HandlerSupport) CreateUser(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "данный http-метод не поддерживается", http.StatusMethodNotAllowed)
		log.Println("был вызван другой тип запроса: ", r.Method)
		return
	}

	var user models.UserRequest
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Ошибка при декодировании объекта", http.StatusBadRequest)
		log.Println("Ошибка при декодировании объекта: ", err)
		return
	}

	err = h.Serv.RegistorUser(user.Email, user.Password)
	if err != nil {
		if strings.Contains(err.Error(), "Email уже занят") {
			http.Error(w, "Email уже зарегистрирован", http.StatusConflict)
			return
		}
		http.Error(w, "Ошибка при регистрации пользователя", http.StatusInternalServerError)
		log.Println("Ошибка при добавлении юзера в бд: ", err)
		return
	}
	response := map[string]string{
		"message": "user registered",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Println("ошибка при кодировки JSON: ", err)
		return
	}

}

func (h *HandlerSupport) LoginUser(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Данный метод не подддерживается", http.StatusMethodNotAllowed)
		return
	}

	var user models.UserRequest
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Ошибка при декодировании данных", http.StatusBadRequest)
		log.Println("Ошибка при декодировании данных: ", err)
		return
	}

	userID, err := h.Serv.CheckLoginUser(user.Email, user.Password)
	if err != nil {
		http.Error(w, "неверные учетные данные", http.StatusUnauthorized)
		return
	}

	token, err := h.Serv.GenerateJWT(userID)
	if err != nil {
		http.Error(w, "внутрення ошибка сервера связанная с токеном", http.StatusInternalServerError)
		return
	}
	response := map[string]string{
		"token": token,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)

}

func (h *HandlerSupport) Me(w http.ResponseWriter, r *http.Request) {

	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		http.Error(w, "пользователь не авторизован", http.StatusUnauthorized)
		return
	}
	log.Printf("Запрос от пользователя с ID: %d", userID)

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(userID)
	if err != nil {
		http.Error(w, "Ошибка при декодировании", http.StatusUnauthorized)
		return
	}
}

func (h *HandlerSupport) CreateContact(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "данный метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}

	var modelCreateContact models.CreateContactRequest
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&modelCreateContact)
	if err != nil {
		http.Error(w, "некорректный формат JSON или неверные типы данных", http.StatusBadRequest)
		return
	}

	err = h.Serv.CreateContact(userID, modelCreateContact.Name, modelCreateContact.Email, modelCreateContact.Phone)
	if err != nil {
		http.Error(w, "ошибка авторизации", http.StatusUnauthorized)
		return
	}
	response := map[string]string{
		"message": "contact created",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Println("ошибка при кодировки JSON: ", err)
		return
	}

}
