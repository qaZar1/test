package service

import (
	"database/sql"
	"errors"
	"net/http"

	jsoniter "github.com/json-iterator/go"
	"github.com/qaZar1/test/wallet/autogen"
	"github.com/qaZar1/test/wallet/internal/postgres"
	"github.com/qaZar1/test/wallet/internal/utils"
)

type Transport struct {
	service *Service
}

func NewTransport(cfg postgres.Config) autogen.ServerInterface {
	return &Transport{
		service: NewService(cfg),
	}
}

// Обновление баланса кошелька
// (POST /api/v1/wallet)
func (t *Transport) PostApiV1Wallet(w http.ResponseWriter, r *http.Request) {
	var wallet autogen.WalletUpdate
	if err := jsoniter.NewDecoder(r.Body).Decode(&wallet); err != nil {
		utils.WriteString(w, http.StatusBadRequest, ErrUnmarshalBody.Error())
		return
	}

	if err := t.service.UpsertWallet(wallet); err != nil {
		switch {
		case errors.Is(err, ErrInvalidOperationType):
			utils.WriteString(w, http.StatusBadRequest, ErrInvalidOperationType.Error())
		case errors.Is(err, ErrNotEnoughFunds):
			utils.WriteString(w, http.StatusUnprocessableEntity, ErrNotEnoughFunds.Error())
		default:
			utils.WriteString(w, http.StatusInternalServerError, ErrInternalError.Error())
		}
		return
	}

	utils.WriteNoContent(w)
}

// Получение баланса кошелька
// (GET /api/v1/wallets/{wallet_id})
func (t *Transport) GetApiV1WalletsWalletId(w http.ResponseWriter, r *http.Request, walletId string) {
	wallet, err := t.service.Get(walletId)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			utils.WriteString(w, http.StatusNotFound, ErrWalletNotFound.Error())
		case errors.Is(err, ErrInvalidID):
			utils.WriteString(w, http.StatusBadRequest, ErrInvalidID.Error())
		case errors.Is(err, ErrInvalidSyntax):
			utils.WriteString(w, http.StatusBadRequest, ErrInvalidSyntax.Error())
		default:
			utils.WriteString(w, http.StatusInternalServerError, ErrGetWallet.Error())
		}
		return
	}

	utils.WriteJSON(w, http.StatusOK, wallet)
}
