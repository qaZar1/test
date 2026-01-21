package postgres

import (
	"database/sql"
	"errors"
	"strconv"

	"github.com/qaZar1/test/wallet/autogen"
)

var (
	ErrNotEnoughFunds = errors.New("not enough funds")
	ErrWalletNotFound = errors.New("wallet not found")
)

func (pg *postgres) GetWallet(walletID string) (*autogen.Wallet, error) {
	const query = `
SELECT wallet_id, amount
FROM wallets.balances
WHERE wallet_id = $1;
	`

	rows := pg.db.QueryRow(query, walletID)
	if rows == nil {
		return nil, sql.ErrNoRows
	}

	var walletIDStr string
	var amountStr string
	if err := rows.Scan(&walletIDStr, &amountStr); err != nil {
		return nil, err
	}

	if walletIDStr == "" {
		return nil, ErrWalletNotFound
	}

	amountInt, err := strconv.Atoi(amountStr)
	if err != nil {
		return nil, err
	}

	return &autogen.Wallet{
		WalletId: walletIDStr,
		Amount:   int64(amountInt),
	}, nil
}

func (pg *postgres) UpsertWallet(wallet *autogen.WalletUpdate) error {
	const query = `
INSERT INTO wallets.balances (wallet_id, amount)
VALUES ($1, $2)
ON CONFLICT (wallet_id) DO UPDATE
SET amount = wallets.balances.amount + $2
WHERE wallets.balances.amount + $2 >= 0;
`

	result, err := pg.db.Exec(query, wallet.WalletId, wallet.Amount)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return ErrNotEnoughFunds
	}

	return nil
}
