package service_test

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/golang/mock/gomock"
	"github.com/qaZar1/test/autogen"
	"github.com/qaZar1/test/autogen/mocks"
	"github.com/qaZar1/test/internal/service"
	"github.com/stretchr/testify/assert"
)

func TestGetBalance(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mocks.NewMockIRepository(ctrl)

		walletID := "6321942b-2ffb-49a6-930e-59fb43c2523a"

		mockRepo.EXPECT().GetWallet(walletID).Return(&autogen.Wallet{
			WalletId: walletID,
			Amount:   100,
		}, nil).Times(1)

		transport := service.Transport{
			Repo:      mockRepo,
			Validator: validator.New(),
		}

		req := httptest.NewRequest(http.MethodGet, "/api/v1/wallets/"+walletID, nil)
		w := httptest.NewRecorder()

		transport.GetApiV1WalletsWalletId(w, req, walletID)

		res := w.Result()

		var resp autogen.Wallet
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, walletID, resp.WalletId)
		assert.Equal(t, int64(100), resp.Amount)

		if res.StatusCode != http.StatusOK {
			t.Errorf("expected status 200, got %d", res.StatusCode)
		}
	})
	t.Run("ErrorInvalidID", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mocks.NewMockIRepository(ctrl)

		walletID := "qwe"

		mockRepo.EXPECT().GetWallet(walletID).Return(nil, service.ErrInvalidID).Times(1)

		transport := service.Transport{
			Repo:      mockRepo,
			Validator: validator.New(),
		}

		req := httptest.NewRequest(http.MethodGet, "/api/v1/wallets/"+walletID, nil)
		w := httptest.NewRecorder()

		transport.GetApiV1WalletsWalletId(w, req, walletID)

		res := w.Result()

		if res.StatusCode != http.StatusBadRequest {
			t.Errorf("expected status 400, got %d", res.StatusCode)
		}
	})

	t.Run("ErrorNotFound", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mocks.NewMockIRepository(ctrl)

		walletID := "6321942b-2ffb-49a6-930e-59fb43c2523f"

		mockRepo.EXPECT().GetWallet(walletID).Return(nil, sql.ErrNoRows).Times(1)

		transport := service.Transport{
			Repo:      mockRepo,
			Validator: validator.New(),
		}

		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/wallets/%s", walletID), nil)
		w := httptest.NewRecorder()

		transport.GetApiV1WalletsWalletId(w, req, walletID)

		res := w.Result()

		if res.StatusCode != http.StatusNotFound {
			t.Errorf("expected status 404, got %d", res.StatusCode)
		}
	})

	t.Run("ErrorInternal", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mocks.NewMockIRepository(ctrl)

		walletID := "6321942"

		mockRepo.EXPECT().GetWallet(walletID).Return(nil, service.ErrInternalError).Times(1)

		transport := service.Transport{
			Repo: mockRepo,
		}

		req := httptest.NewRequest(http.MethodGet, "/api/v1/wallets/"+walletID, nil)
		w := httptest.NewRecorder()

		transport.GetApiV1WalletsWalletId(w, req, walletID)

		res := w.Result()

		if res.StatusCode != http.StatusInternalServerError {
			t.Errorf("expected status 500, got %d", res.StatusCode)
		}
	})
}

func TestPostWalletDeposit(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mocks.NewMockIRepository(ctrl)

		body := `
	{
		"wallet_id":"6321942b-2ffb-49a6-930e-59fb43c2523a",
		"operation_type":"DEPOSIT",
		"amount":100
	}`

		mockRepo.EXPECT().UpsertWallet(autogen.WalletUpdate{
			WalletId:      "6321942b-2ffb-49a6-930e-59fb43c2523a",
			OperationType: "DEPOSIT",
			Amount:        100,
		}).Times(1)

		transport := service.Transport{
			Repo:      mockRepo,
			Validator: validator.New(),
		}

		req := httptest.NewRequest(
			http.MethodPost,
			"/api/v1/wallet",
			strings.NewReader(body),
		)
		w := httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")

		transport.PostApiV1Wallet(w, req)

		res := w.Result()

		if res.StatusCode != http.StatusNoContent {
			t.Errorf("expected status 204, got %d", res.StatusCode)
		}
	})
	t.Run("ErrorValidation", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mocks.NewMockIRepository(ctrl)

		walletUpdate := `
	{
		"wallet_id": "6321942b-2ffb-49a6-930e-59fb43c2523a",
		"operation_type": "WITHDRAW",
		"amount": 1000
	}`

		mockRepo.EXPECT().
			UpsertWallet(
				autogen.WalletUpdate{
					WalletId:      "6321942b-2ffb-49a6-930e-59fb43c2523a",
					OperationType: "WITHDRAW",
					Amount:        1000,
				},
			).Times(1)

		transport := service.Transport{
			Repo:      mockRepo,
			Validator: validator.New(),
		}

		req := httptest.NewRequest(
			http.MethodPost,
			"/api/v1/wallet",
			strings.NewReader(walletUpdate),
		)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		transport.PostApiV1Wallet(w, req)

		res := w.Result()

		if res.StatusCode != http.StatusNoContent {
			t.Errorf("expected status 204, got %d", res.StatusCode)
		}
	})
	t.Run("ErrorValidation", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mocks.NewMockIRepository(ctrl)

		body := `{
		"wallet_id": "6321942b-2ffb-49a6-930e-59fb43c2523a",
		"amount": 100,
		"operation_type": "DEPOSIT
	}`

		transport := service.Transport{
			Repo: mockRepo,
		}

		req := httptest.NewRequest(
			http.MethodPost,
			"/api/v1/wallet",
			strings.NewReader(body),
		)
		w := httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")

		transport.PostApiV1Wallet(w, req)

		res := w.Result()

		if res.StatusCode != http.StatusBadRequest {
			t.Errorf("expected status 400, got %d", res.StatusCode)
		}
	})

	t.Run("ErrorValidationFailed", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mocks.NewMockIRepository(ctrl)

		body := `
	{
		"wallet_id": "6321942b-2ffb-49a6-930e-59fb43c2523a",
		"operation_type": "DEBIT",
		"amount": 1000
	}`

		transport := service.Transport{
			Repo:      mockRepo,
			Validator: validator.New(),
		}

		req := httptest.NewRequest(
			http.MethodPost,
			"/api/v1/wallet",
			strings.NewReader(body),
		)

		w := httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")

		transport.PostApiV1Wallet(w, req)

		res := w.Result()

		if res.StatusCode != http.StatusBadRequest {
			t.Errorf("expected status 400, got %d", res.StatusCode)
		}

		respBody, _ := io.ReadAll(res.Body)

		assert.Contains(t, string(respBody), service.ErrValidationFailed.Error())
	})

	t.Run("ErrorNotEnoughFunds", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mocks.NewMockIRepository(ctrl)

		body := `
	{
		"wallet_id": "6321942b-2ffb-49a6-930e-59fb43c2523a",
		"operation_type": "WITHDRAW",
		"amount": 1000
	}`

		mockRepo.EXPECT().
			UpsertWallet(
				autogen.WalletUpdate{
					WalletId:      "6321942b-2ffb-49a6-930e-59fb43c2523a",
					OperationType: "WITHDRAW",
					Amount:        1000,
				},
			).Return(service.ErrNotEnoughFunds).Times(1)

		transport := service.Transport{
			Repo:      mockRepo,
			Validator: validator.New(),
		}

		req := httptest.NewRequest(
			http.MethodPost,
			"/api/v1/wallet",
			strings.NewReader(body),
		)

		w := httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")

		transport.PostApiV1Wallet(w, req)

		res := w.Result()

		if res.StatusCode != http.StatusUnprocessableEntity {
			t.Errorf("expected status 422, got %d", res.StatusCode)
		}
	})

	t.Run("ErrorInternal", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mocks.NewMockIRepository(ctrl)

		body := `
	{
		"wallet_id": "6321942b-2ffb-49a6-930e-59fb43c2523a",
		"operation_type": "WITHDRAW",
		"amount": 1000
	}`

		mockRepo.EXPECT().
			UpsertWallet(
				autogen.WalletUpdate{
					WalletId:      "6321942b-2ffb-49a6-930e-59fb43c2523a",
					OperationType: "WITHDRAW",
					Amount:        1000,
				},
			).Return(service.ErrInternalError).Times(1)

		transport := service.Transport{
			Repo:      mockRepo,
			Validator: validator.New(),
		}

		req := httptest.NewRequest(
			http.MethodPost,
			"/api/v1/wallet",
			strings.NewReader(body),
		)

		w := httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")

		transport.PostApiV1Wallet(w, req)

		res := w.Result()

		if res.StatusCode != http.StatusInternalServerError {
			t.Errorf("expected status 500, got %d", res.StatusCode)
		}
	})

	t.Run("WithNegativeAmount", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mocks.NewMockIRepository(ctrl)

		body := `
	{
		"wallet_id": "6321942b-2ffb-49a6-930e-59fb43c2523a",
		"operation_type": "WITHDRAW",
		"amount": -1000
	}`

		transport := service.Transport{
			Repo:      mockRepo,
			Validator: validator.New(),
		}

		req := httptest.NewRequest(
			http.MethodPost,
			"/api/v1/wallet",
			strings.NewReader(body),
		)

		w := httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")

		transport.PostApiV1Wallet(w, req)

		res := w.Result()

		if res.StatusCode != http.StatusBadRequest {
			t.Errorf("expected status 500, got %d", res.StatusCode)
		}

		bodyStr, _ := io.ReadAll(res.Body)

		assert.Equal(t, string(bodyStr), service.ErrValidationFailed.Error())
	})
}
