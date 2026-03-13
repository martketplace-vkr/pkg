package account

import (
	"errors"
)

var (
	ErrTooBigNumberToGenerateAccountID = errors.New("too big number to generate account/wallet id")
	ErrInvalidAccountID                = errors.New("invalid account/wallet id")
	ErrInvalidCurrencyID               = errors.New("invalid currency id")
	ErrNilPtr                          = errors.New("nil pointer")
)
