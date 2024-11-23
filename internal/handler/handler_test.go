package handler_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"wallet-app/internal/handler"
	"wallet-app/internal/service"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestWalletOperations(t *testing.T) {
	mockService := &service.MockWalletService{
		ProcessOperationFn: func(walletID uuid.UUID, operationType string, amount int64) error {
			// Имитация успешной работы
			if operationType == "DEPOSIT" && amount > 0 {
				return nil
			}
			if operationType == "WITHDRAW" && amount > 0 {
				return nil
			}
			if amount <= 0 {
				return fmt.Errorf("amount must be greater than zero")
			}
			return fmt.Errorf("invalid operation type")
		},
	}

	h := handler.New(mockService)

	tests := []struct {
		name           string
		operation      string
		requestBody    map[string]interface{}
		expectedStatus int
		expectedBody   string
	}{
		{
			name:      "Successful DEPOSIT",
			operation: "DEPOSIT",
			requestBody: map[string]interface{}{
				"walletID":      uuid.New(),
				"operationType": "DEPOSIT",
				"amount":        1000,
			},
			expectedStatus: http.StatusOK,
			expectedBody:   "Operation successful",
		},
		{
			name:      "Successful WITHDRAW",
			operation: "WITHDRAW",
			requestBody: map[string]interface{}{
				"walletID":      uuid.New(),
				"operationType": "WITHDRAW",
				"amount":        500,
			},
			expectedStatus: http.StatusOK,
			expectedBody:   "Operation successful",
		},
		{
			name:      "Invalid operation type",
			operation: "INVALID",
			requestBody: map[string]interface{}{
				"walletID":      uuid.New(),
				"operationType": "INVALID",
				"amount":        500,
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "invalid operation type",
		},
		{
			name:      "Negative amount",
			operation: "WITHDRAW",
			requestBody: map[string]interface{}{
				"walletID":      uuid.New(),
				"operationType": "WITHDRAW",
				"amount":        -100,
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "amount must be greater than zero", // надо возвращать более конкретную ошибку
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			r := chi.NewRouter()
			r.Post("/api/v1/wallets/operation", h.HandleWalletOperation)

			server := httptest.NewServer(r)
			defer server.Close()

			reqBody, _ := json.Marshal(test.requestBody)
			resp, err := http.Post(server.URL+"/api/v1/wallets/operation", "application/json", bytes.NewBuffer(reqBody))

			assert.NoError(t, err)
			assert.Equal(t, test.expectedStatus, resp.StatusCode)

			body, _ := io.ReadAll(resp.Body)
			assert.Contains(t, string(body), test.expectedBody)
		})
	}
}
