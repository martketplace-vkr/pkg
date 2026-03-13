package account

import (
	"database/sql"
	"database/sql/driver"
	"encoding"
	"fmt"
	"strings"

	"github.com/martketplace-vkr/pkg/utils/currency"

	"github.com/goccy/go-json"
)

var (
	_ = fmt.Stringer(&Account{})
	_ = encoding.TextMarshaler(&Account{})
	_ = encoding.TextUnmarshaler(&Account{})
	_ = json.Marshaler(&Account{})
	_ = json.Unmarshaler(&Account{})
	_ = sql.Scanner(&Account{})
	_ = driver.Valuer(&Account{})
)

const AccountLength = 30

type Account struct {
	EntityTypeID        int64             // 2 цифры
	EntityAccountTypeID int64             // 2
	CurrencyID          currency.Currency // 4
	CountryID           int64             // 3
	CityID              int64             // 3
	EntityID            int64             // 12
	Number              int64             // 4
}

func GetFromID(accountID string) (a Account, err error) {
	if len(accountID) != AccountLength {
		return a, ErrInvalidAccountID
	}

	destCurrID := int64(a.CurrencyID)

	err = splitAndParse(accountID, []destAndLen{
		{
			dest: &a.EntityTypeID,
			len:  EntityTypeIDLen,
		},
		{
			dest: &a.EntityAccountTypeID,
			len:  EntityAccountTypeIDLen,
		},
		{
			dest: &destCurrID,
			len:  CurrencyIDLen,
		},
		{
			dest: &a.CountryID,
			len:  CountryIDLen,
		},
		{
			dest: &a.CityID,
			len:  CityIDLen,
		},
		{
			dest: &a.EntityID,
			len:  EntityIDLen,
		},
		{
			dest: &a.Number,
			len:  NumberLen,
		},
	})
	if err != nil {
		return a, err
	}

	a.CurrencyID = currency.Currency(destCurrID)

	return a, nil
}

const (
	EntityTypeIDLen        = 2
	EntityAccountTypeIDLen = 2
	CurrencyIDLen          = 4
	CountryIDLen           = 3
	CityIDLen              = 3
	EntityIDLen            = 12
	NumberLen              = 4
)

func (a Account) GetID() (accountID string, err error) {
	b := strings.Builder{}
	b.Grow(AccountLength)

	str := ""
	for _, item := range []numAndLen{
		{
			num: a.EntityTypeID,
			len: EntityTypeIDLen,
		},
		{
			num: a.EntityAccountTypeID,
			len: EntityAccountTypeIDLen,
		},
		{
			num: int64(a.CurrencyID),
			len: CurrencyIDLen,
		},
		{
			num: a.CountryID,
			len: CountryIDLen,
		},
		{
			num: a.CityID,
			len: CityIDLen,
		},
		{
			num: a.EntityID,
			len: EntityIDLen,
		},
		{
			num: a.Number,
			len: NumberLen,
		},
	} {
		str, err = toStringWithLen(item)
		if err != nil {
			return "", err
		}

		b.WriteString(str)
	}

	return b.String(), nil
}
func (a Account) GetWalletID() (accountID string, err error) {
	if a.CountryID == WalletCountryID && a.CityID == WalletCityID {
		return a.GetID()
	}

	wallet := Account{
		EntityTypeID:        a.EntityTypeID,
		EntityAccountTypeID: a.EntityAccountTypeID,
		CurrencyID:          a.CurrencyID,
		CountryID:           WalletCountryID,
		CityID:              WalletCityID,
		EntityID:            a.EntityID,
		Number:              a.Number,
	}

	return wallet.GetID()
}
func (a Account) IsWallet() bool {
	if a.CountryID == WalletCountryID && a.CityID == WalletCityID {
		return true
	}

	return false
}

func (a *Account) String() string {
	accountID, err := a.GetID()
	if err != nil {
		return fmt.Sprintf("%#+v", accountID)
	}

	return accountID
}

func (a Account) MarshalText() (text []byte, err error) {
	accountID, err := a.GetID()
	if err != nil {
		return nil, err
	}

	return []byte(accountID), nil
}
func (a *Account) UnmarshalText(text []byte) error {
	account, err := GetFromID(string(text))
	if err != nil {
		return err
	}

	a.EntityTypeID = account.EntityTypeID
	a.EntityAccountTypeID = account.EntityAccountTypeID
	a.CurrencyID = account.CurrencyID
	a.CountryID = account.CountryID
	a.CityID = account.CityID
	a.EntityID = account.EntityID
	a.Number = account.Number

	return nil
}

func (a Account) MarshalJSON() ([]byte, error) {
	text, err := a.MarshalText()
	if err != nil {
		return nil, err
	}

	return json.Marshal(string(text))
}
func (a *Account) UnmarshalJSON(data []byte) error {
	text := ""

	err := json.Unmarshal(data, &text)
	if err != nil {
		return err
	}

	return a.UnmarshalText([]byte(text))
}

func (a *Account) Scan(src any) error {
	if v, ok := src.(string); ok {
		return a.UnmarshalText([]byte(v))
	}
	if v, ok := src.(*string); ok {
		if v == nil {
			return ErrNilPtr
		}

		return a.UnmarshalText([]byte(*v))
	}
	if v, ok := src.([]byte); ok {
		return a.UnmarshalText(v)
	}
	if v, ok := src.(*[]byte); ok {
		if v == nil {
			return ErrNilPtr
		}

		return a.UnmarshalText(*v)
	}

	return ErrInvalidAccountID
}

func (a Account) Value() (driver.Value, error) {
	accountID, err := a.GetID()
	if err != nil {
		return nil, err
	}

	return driver.Value(accountID), nil
}

func (a Account) IsCryptoWallet() (bool, error) {
	return CurrencyIsCrypto(a.CurrencyID)
}

func (a Account) IsFiatWallet() (bool, error) {
	return CurrencyIsFiat(a.CurrencyID)
}
