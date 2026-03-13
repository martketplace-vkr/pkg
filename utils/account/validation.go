package account

import (
	"reflect"

	"github.com/go-playground/validator/v10"
)

var (
	accountType    = reflect.TypeOf(Account{})
	accountPtrType = reflect.TypeOf(&Account{})
)

func RegisterValidation(v *validator.Validate) (err error) {
	err = v.RegisterValidation("account", AccountValidation)
	if err != nil {
		return err
	}

	err = v.RegisterValidation("wallet", WalletValidation)
	if err != nil {
		return err
	}

	err = v.RegisterValidation("fiatAccount", FiatAccountValidation)
	if err != nil {
		return err
	}

	err = v.RegisterValidation("clientAccount", ClientAccountValidation)
	if err != nil {
		return err
	}

	return nil
}

func basicValidation(fl validator.FieldLevel) (accountID string, account Account, valid bool) {
	// Общая валидация
	ok := false
	err := error(nil)

	k := fl.Field().Kind()
	if k == reflect.String {
		// Если валидируем строку (id счёта)
		accountID = fl.Field().String()

		account, err = GetFromID(accountID)
		if err != nil {
			return accountID, account, false
		}
	} else if k == reflect.Struct && fl.Field().CanConvert(accountType) {
		account, ok = fl.Field().Convert(accountType).Interface().(Account)
		if !ok {
			return accountID, account, false
		}

		accountID, err = account.GetID()
		if err != nil {
			return accountID, account, false
		}
	} else if k == reflect.Struct && fl.Field().CanConvert(accountPtrType) {
		accountPtr := (*Account)(nil)

		accountPtr, ok = fl.Field().Convert(accountPtrType).Interface().(*Account)
		if !ok {
			return accountID, account, false
		}

		account = *accountPtr

		accountID, err = account.GetID()
		if err != nil {
			return accountID, account, false
		}
	}

	// Валидируем валюту

	isCrypto, err := CurrencyIsCrypto(account.CurrencyID)
	if err != nil {
		return accountID, account, false
	}

	isFiat, err := CurrencyIsFiat(account.CurrencyID)
	if err != nil {
		return accountID, account, false
	}

	if !isCrypto && !isFiat {
		return accountID, account, false
	}

	// Либо страна и город оба пустые (для кошельков), либо оба заполнены (для счетов)
	if (account.CountryID == WalletCountryID) != (account.CityID == WalletCityID) {
		return accountID, account, false
	}

	return accountID, account, true
}

func AccountValidation(fl validator.FieldLevel) bool {
	_, _, valid := basicValidation(fl)
	
	return valid
}
func WalletValidation(fl validator.FieldLevel) bool {
	_, account, valid := basicValidation(fl)
	if !valid {
		return false
	}

	if account.CountryID != 0 || account.CityID != 0 {
		return false
	}

	return true
}
func FiatAccountValidation(fl validator.FieldLevel) bool {
	_, account, valid := basicValidation(fl)
	if !valid {
		return false
	}

	isFiat, err := CurrencyIsFiat(account.CurrencyID)
	if err != nil {
		return false
	}

	if !isFiat {
		return false
	}

	// У фиатного счёта обязательно должен быть город
	if account.CountryID == 0 || account.CityID == 0 {
		return false
	}

	return true
}
func ClientAccountValidation(fl validator.FieldLevel) bool {
	_, account, valid := basicValidation(fl)
	if !valid {
		return false
	}

	if account.EntityTypeID != EntityClient {
		return false
	}

	return true
}
