package service

import (
	"strings"

	"github.com/google/uuid"
	"github.com/qaZar1/test/wallet/autogen"
	"github.com/qaZar1/test/wallet/internal/postgres"
)

type Service struct {
	db postgres.IPostgres
}

func NewService(cfg postgres.Config) *Service {
	return &Service{
		db: postgres.NewPostgres(cfg),
	}
}

func (s *Service) UpsertWallet(wallet autogen.WalletUpdate) error {
	operationType := strings.ToLower(wallet.OperationType)
	if operationType != DepositOperationType && operationType != WithdrawOperationType {
		return ErrInvalidOperationType
	}

	if operationType == WithdrawOperationType {
		wallet.Amount = -wallet.Amount
	}

	return s.db.UpsertWallet(&wallet)
}

func (s *Service) Get(walletId string) (*autogen.Wallet, error) {
	if err := uuid.Validate(walletId); err != nil {
		return nil, ErrInvalidID
	}

	return s.db.GetWallet(walletId)
}
