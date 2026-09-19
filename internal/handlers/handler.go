package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/yukay/CRM/internal/models"
	"github.com/yukay/CRM/internal/service"
)

type HandlerSupport struct {
	Serv *service.ServiceSupport
}

func NewHandlerSupport(serv *service.ServiceSupport) *HandlerSupport {
	return &HandlerSupport{Serv: serv}
}

func Ping(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	w.Write([]byte(`{"message": "pong"}`))
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

	w.WriteHeader(http.StatusCreated)

}
