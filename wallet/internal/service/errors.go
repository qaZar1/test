package service

import "errors"

var (
	ErrUnmarshalBody        = errors.New("failed to unmarshal request body")
	ErrInvalidOperationType = errors.New("invalid operation type")
	ErrInvalidRequestBody   = errors.New("invalid request body")
	ErrNotEnoughFunds       = errors.New("not enough funds")
	ErrInternalError        = errors.New("internal error")

	ErrInvalidID      = errors.New("invalid wallet ID")
	ErrWalletNotFound = errors.New("wallet not found")
	ErrGetWallet      = errors.New("failed to get wallet")
	ErrInvalidSyntax  = errors.New("invalid syntax")
)
