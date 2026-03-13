package account

import (
	"strconv"
	"strings"

	"github.com/martketplace-vkr/pkg/utils/countries"
	"github.com/martketplace-vkr/pkg/utils/currency"
)

func CurrencyIsCrypto(currencyID currency.Currency) (isCrypto bool, err error) {
	str, err := toStringWithLen(numAndLen{
		num: int64(currencyID),
		len: CurrencyIDLen,
	})
	if err != nil {
		return false, ErrInvalidCurrencyID
	}

	isCrypto = str[0] == '2'

	return isCrypto, nil
}
func CurrencyIsFiat(currencyID currency.Currency) (isFiat bool, err error) {
	str, err := toStringWithLen(numAndLen{
		num: int64(currencyID),
		len: CurrencyIDLen,
	})
	if err != nil {
		return false, ErrInvalidCurrencyID
	}

	isFiat = str[0] == '1'

	return isFiat, nil
}

func GetClientWallet(currencyID currency.Currency, clientID, number int64) Account {
	return Account{
		EntityTypeID:        EntityClient,
		EntityAccountTypeID: AccountTypeClientBalance,
		CurrencyID:          currencyID,
		CountryID:           WalletCountryID,
		CityID:              WalletCityID,
		EntityID:            clientID,
		Number:              number,
	}
}

// Эти штуки нужны для парсинга/форматирования id счёта

type numAndLen struct {
	num int64
	len int
}

func toStringWithLen(item numAndLen) (string, error) {
	numStr := strconv.FormatInt(item.num, 10)
	if len(numStr) > item.len {
		return "", ErrTooBigNumberToGenerateAccountID
	}

	if zeroesCount := item.len - len(numStr); zeroesCount > 0 {
		numStr = strings.Repeat("0", zeroesCount) + numStr
	}

	return numStr, nil
}

type destAndLen struct {
	dest *int64
	len  int
}

func splitAndParse(str string, items []destAndLen) (err error) {
	start := 0
	for i := range items {
		err = fromStringToDest(strAndDest{
			str:  str[start : start+items[i].len],
			dest: items[i].dest,
		})
		if err != nil {
			return err
		}

		start += items[i].len
	}

	return nil
}

type strAndDest struct {
	str  string
	dest *int64
}

func fromStringToDest(item strAndDest) error {
	if item.dest == nil {
		return ErrInvalidAccountID
	}

	err := error(nil)

	*item.dest, err = strconv.ParseInt(item.str, 10, 64)
	if err != nil {
		return ErrInvalidAccountID
	}

	return nil
}

func GetClientReferralAccounts(clientID int64) (active, freeze Account) {
	return Account{
			EntityTypeID:        EntityClient,
			EntityAccountTypeID: AccountTypeClientReferralBalance,
			CurrencyID:          currency.USDTinTRC,
			CountryID:           int64(countries.CountryRUID),
			CityID:              int64(countries.CityRUMoscowID),
			EntityID:            clientID,
			Number:              1,
		}, Account{
			EntityTypeID:        EntityClient,
			EntityAccountTypeID: AccountTypeClientReferralFreezeBalance,
			CurrencyID:          currency.USDTinTRC,
			CountryID:           int64(countries.CountryRUID),
			CityID:              int64(countries.CityRUMoscowID),
			EntityID:            clientID,
			Number:              1,
		}
}

func GetEngineReferralAccounts() (active, freeze Account) {
	return Account{
			EntityTypeID:        EntityEngine,
			EntityAccountTypeID: AccountTypeReferralEngine,
			CurrencyID:          currency.USDTinTRC,
			CountryID:           int64(countries.CountryRUID),
			CityID:              int64(countries.CityRUMoscowID),
			EntityID:            0,
			Number:              0,
		}, Account{
			EntityTypeID:        EntityEngine,
			EntityAccountTypeID: AccountTypeReferralEngineFreeze,
			CurrencyID:          currency.USDTinTRC,
			CountryID:           int64(countries.CountryRUID),
			CityID:              int64(countries.CityRUMoscowID),
			EntityID:            0,
			Number:              0,
		}
}
