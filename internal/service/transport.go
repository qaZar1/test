package service

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"
	jsoniter "github.com/json-iterator/go"
	"github.com/qaZar1/test/autogen"
	"github.com/qaZar1/test/internal/repository"
)

type Transport struct {
	Repo      repository.IRepository
	Validator *validator.Validate
}

func NewTransport(repo repository.IRepository) *Transport {
	return &Transport{
		Repo:      repo,
		Validator: validator.New(),
	}
}

// Обновление баланса кошелька
// (POST /api/v1/wallet)
func (t *Transport) PostApiV1Wallet(w http.ResponseWriter, r *http.Request) {
	var wallet autogen.WalletUpdate
	if err := jsoniter.NewDecoder(r.Body).Decode(&wallet); err != nil {
		writeString(w, http.StatusBadRequest, ErrUnmarshalBody.Error())
		return
	}

	if err := t.Validator.Struct(wallet); err != nil {
		writeString(w, http.StatusBadRequest, ErrValidationFailed.Error())
		return
	}

	if err := t.Repo.UpsertWallet(wallet); err != nil {
		switch {
		case errors.Is(err, ErrNotEnoughFunds):
			writeString(w, http.StatusUnprocessableEntity, ErrNotEnoughFunds.Error())
		default:
			writeString(w, http.StatusInternalServerError, ErrInternalError.Error())
		}
		return
	}

	writeNoContent(w)
}

// Получение баланса кошелька
// (GET /api/v1/wallets/{wallet_id})
func (t *Transport) GetApiV1WalletsWalletId(w http.ResponseWriter, r *http.Request, walletId string) {
	wallet, err := t.Repo.GetWallet(walletId)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			writeString(w, http.StatusNotFound, ErrWalletNotFound.Error())
		case errors.Is(err, ErrInvalidID):
			writeString(w, http.StatusBadRequest, ErrInvalidID.Error())
		default:
			writeString(w, http.StatusInternalServerError, ErrGetWallet.Error())
		}
		return
	}

	writeJSON(w, http.StatusOK, wallet)
}
