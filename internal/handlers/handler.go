package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
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

func (h *HandlerSupport) GetContacts(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "данный метод недоступен", http.StatusMethodNotAllowed)
		return
	}

	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		http.Error(w, "не авторизован", http.StatusUnauthorized)
		return
	}

	contacts, err := h.Serv.GetContacts(userID)
	if err != nil {
		http.Error(w, "внутрення ошибка сервера", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(contacts)
	if err != nil {
		log.Println("ошибка при кодировании данных: ", err)
		return
	}

}

func (h *HandlerSupport) GetContactByID(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "данный метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		http.Error(w, "ошибка авторизации", http.StatusUnauthorized)
		return
	}

	idStr := r.PathValue("id")
	contactID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ошибка невалидных данных", http.StatusBadRequest)
		return
	}
	contact, err := h.Serv.GetContactByID(userID, contactID)
	if err != nil {
		if err.Error() == "нет контакта в бд" {
			http.Error(w, "сущность не найдена в бд", http.StatusNotFound)
			return
		}

		log.Printf("ошибка получения контакта (userID: %d, contactID: %d): %v", userID, contactID, err)
		http.Error(w, "внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(contact)
	if err != nil {
		log.Println("ошибка при кодировании: ", err)
		return
	}

}

func (h *HandlerSupport) DeleteContactByID(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodDelete {
		http.Error(w, "данный метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		http.Error(w, "ошибка авторизации", http.StatusUnauthorized)
		return
	}

	idStr := r.PathValue("id")
	contactID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ошибка невалидных данных", http.StatusBadRequest)
		return
	}

	err = h.Serv.DeleteContactByID(userID, contactID)
	if err != nil {
		if err.Error() == "not found" {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}

		http.Error(w, "внутренняя ошибка сервера", http.StatusInternalServerError)
		log.Println("ошибка при удалении контакта: ", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)

}

func (h *HandlerSupport) UpdateContact(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPut {
		http.Error(w, "данный метод не поддерживаетя", http.StatusMethodNotAllowed)
		return
	}

	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		http.Error(w, "ошибка авторизации", http.StatusUnauthorized)
		return
	}

	idStr := r.PathValue("id")
	contactID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ошибка валидации данных", http.StatusBadRequest)
		return
	}

	var request models.UpdateContact
	defer r.Body.Close()
	err = json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "ошибка JSON формата", http.StatusBadRequest)
		return
	}

	err = h.Serv.UpdateContact(userID, contactID, request.Name, request.Email, request.Phone)
	if err != nil {
		if err.Error() == "not found" {
			http.Error(w, "объект не найден", http.StatusNotFound)
			return
		}

		http.Error(w, "внутренняя ошибка сервера", http.StatusInternalServerError)
		log.Println("ошибка при обновлении объекта: ", err)
		return
	}

	w.WriteHeader(http.StatusOK)

}
