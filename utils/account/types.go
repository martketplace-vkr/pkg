package account

// Entities and AccountsTypes
//
//goland:noinspection ALL
const (
	EntitySystem                  = 10
	EntityOffice                  = 20
	EntityEngine                  = 30
	EntityClient                  = 40
	EntityOtherActivesAndPassives = 50
	EntityFinancialResults        = 70
	EntityExternalAccounts        = 99

	AccountTypeSystemAuthorizedCapital = 20

	AccountTypeOfficeBalance       = 00
	AccountTypeOfficeFreezeBalance = 01

	AccountTypeEngineBalance        = 00
	AccountTypeEngineFreezeBalance  = 01
	AccountTypeReferralEngine       = 10
	AccountTypeReferralEngineFreeze = 11

	AccountTypeClientBalance               = 00
	AccountTypeClientFreezeBalance         = 01
	AccountTypeClientReferralBalance       = 10
	AccountTypeClientReferralFreezeBalance = 11

	AccountTypeOtherActivesAndPassivesExchangesSettlementsPassive = 00
	AccountTypeOtherActivesAndPassivesExchangesSettlementsActive  = 01

	AccountTypeFinancialResultsIncomeFromTransfers            = 00
	AccountTypeFinancialResultsIncomeFromConversion           = 01
	AccountTypeFinancialResultsOtherIncomes                   = 02
	AccountTypeFinancialResultsIncomeFromOperationsWithCash   = 03
	AccountTypeFinancialResultsIncomeFromOperationsWithFiat   = 04
	AccountTypeFinancialResultsIncomeFromOperationsWithCrypto = 05

	AccountTypeFinancialResultsExpensesForTransfers            = 10
	AccountTypeFinancialResultsExpensesForConversion           = 11
	AccountTypeFinancialResultsOtherExpenses                   = 12
	AccountTypeFinancialResultsExpensesForOperationsWithCash   = 13
	AccountTypeFinancialResultsExpensesForOperationsWithFiat   = 14
	AccountTypeFinancialResultsExpensesForOperationsWithCrypto = 15

	AccountTypeFinancialResultsSummaryMargin = 99

	AccountTypeExternalAccount = 99
)

// Entities IDs and Accounts Numbers
//
//goland:noinspection ALL
const (
	EntityIDCryptoEngine      = 000000000000
	AccountNumberCryptoEngine = 0000
	EntityIDFiatEngine        = 000000000001
	AccountNumberFiatEngine   = 0000
)

const (
	WalletCountryID = int64(0)
	WalletCityID    = int64(0)

	CryptoCountryID = int64(1)
	CryptoCityID    = int64(1)
)
