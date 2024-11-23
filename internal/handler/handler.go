package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"wallet-app/internal/service"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type Handler interface {
	HandleWalletOperation(w http.ResponseWriter, r *http.Request)
	GetWalletBalance(w http.ResponseWriter, r *http.Request)
}

// toDo: comment this
type handlerImpl struct {
	service service.WalletService
}

// toDo: comment this
func New(srv service.WalletService) *handlerImpl {
	return &handlerImpl{service: srv}
}

// toDo: comment this
type WalletOperationRequest struct {
	WalletID      uuid.UUID `json:"walletID"`
	OperationType string    `json:"operationType"`
	Amount        int64     `json:"amount"`
}

// POST запрос -- отправляем данные на /api/v1/wallet
func (h *handlerImpl) HandleWalletOperation(w http.ResponseWriter, r *http.Request) {
	var err error
	// Создаём объект запроса
	var request WalletOperationRequest

	fmt.Printf("Request Path Wallet: %s\n", r.URL.Path)

	// Декодируем JSON в запросе
	if err = json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Выполняем операцию с кошельком, получив все параметры кошелька
	err = h.service.ProcessOperation(request.WalletID, request.OperationType, request.Amount)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// В случае отсутствия ошибок возвращаем статус ОК
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Operation successfully completed"))

}

func (h *handlerImpl) GetWalletBalance(w http.ResponseWriter, r *http.Request) {
	// Получаем id кошелька
	walletIDStr := chi.URLParam(r, "walletID")

	walletID, err := uuid.Parse(walletIDStr)
	if err != nil {
		// if walletID.String() == "" {
		// 	http.Error(w, "wallet ID is required", http.StatusBadRequest)
		// 	return
		// }
		http.Error(w, "invalid UUID", http.StatusBadRequest)
		return
	}

	// Выполняем операцию получения баланса кошелька
	balance, err := h.service.GetBalance(walletID)
	if err != nil {
		if errors.Is(err, service.ErrWalletNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Сохраняем баланс в мапу и возвращаем результат
	resp := map[string]int64{"balance": balance}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
