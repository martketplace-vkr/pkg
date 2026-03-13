package currency

type Currency int64
type CurrencyName string

const (
	// Fiat.
	RUB Currency = 1000
	USD Currency = 1010
	EUR Currency = 1020
	AED Currency = 1030
	UZS Currency = 1040
	KZT Currency = 1050
	CNY Currency = 1060
	GBP Currency = 1070
	TJS Currency = 1080
	KGS Currency = 1090
	GEL Currency = 1100
	MDL Currency = 1110
	AMD Currency = 1120

	// Crypto.
	USDTinERC Currency = 2000
	USDTinTRC Currency = 2001
	USDCinERC Currency = 2002
	USDCinTRC Currency = 2003
	BTC       Currency = 2100
	ETH       Currency = 2110
	TRX       Currency = 2120
	TON       Currency = 2130
)

func IsFiat(currency Currency) bool {
	return currency >= 1000 && currency < 2000
}

func IsCrypto(currency Currency) bool {
	return currency >= 2000 && currency < 3000
}

const (
	// Fiat.
	RUBName CurrencyName = "RUB" // Russian Ruble
	USDName CurrencyName = "USD" // United States Dollar
	EURName CurrencyName = "EUR" // Euro
	AEDName CurrencyName = "AED" // United Arab Emirates Dirham
	UZSName CurrencyName = "UZS" // Uzbek Som
	KZTName CurrencyName = "KZT" // Kazakhstani Tenge
	CNYName CurrencyName = "CNY" // Chinese Yuan
	GBPName CurrencyName = "GBP" // British Pound
	TJSName CurrencyName = "TJS" // Tajikistani Somoni
	KGSName CurrencyName = "KGS" // Kyrgyzstani Som
	GELName CurrencyName = "GEL" // Georgian Lari
	MDLName CurrencyName = "MDL" // Moldovan Leu
	AMDName CurrencyName = "AMD" // Armenian Dram

	// Crypto.
	USDTinERCName CurrencyName = "USDTinERC" // Tether in ERC20
	USDTinTRCName CurrencyName = "USDTinTRC" // Tether in TRC20
	USDCinERCName CurrencyName = "USDCinERC" // USD Coin in ERC20
	USDCinTRCName CurrencyName = "USDCinTRC"
	BTCName       CurrencyName = "BTC" // Bitcoin
	ETHName       CurrencyName = "ETH" // Ethereum
	TRXName       CurrencyName = "TRX" // TRON
	TONName       CurrencyName = "TON" // Toncoin
)

var ExchangedCurrencies = map[Currency][]Currency{
	RUB: {
		USDTinERC,
		USDTinTRC,
		USDCinERC,
		BTC,
		ETH,
		TRX,
		TON,
	},
	USD: {
		USDTinERC,
		USDTinTRC,
		USDCinERC,
		BTC,
		ETH,
		TRX,
		TON,
	},
	EUR: {
		USDTinERC,
		USDTinTRC,
		USDCinERC,
		BTC,
		ETH,
		TRX,
		TON,
	},
	AED: {
		USDTinERC,
		USDTinTRC,
		USDCinERC,
		BTC,
		ETH,
		TRX,
		TON,
	},
	USDTinERC: {
		RUB,
		USD,
		EUR,
		AED,
		BTC,
		ETH,
		TRX,
		TON,
		USDTinTRC,
		USDCinERC,
	},
	USDTinTRC: {
		RUB,
		USD,
		EUR,
		AED,
		BTC,
		ETH,
		TRX,
		TON,
		USDTinERC,
		USDCinERC,
	},
	USDCinERC: {
		RUB,
		USD,
		EUR,
		AED,
		BTC,
		ETH,
		TRX,
		TON,
		USDTinERC,
		USDTinTRC,
	},
	BTC: {
		RUB,
		USD,
		EUR,
		AED,
		USDTinERC,
		USDTinTRC,
		USDCinERC,
	},
	ETH: {
		RUB,
		USD,
		EUR,
		AED,
		USDTinERC,
		USDTinTRC,
		USDCinERC,
	},
	TRX: {
		RUB,
		USD,
		EUR,
		AED,
		USDTinERC,
		USDTinTRC,
		USDCinERC,
	},
	TON: {
		RUB,
		USD,
		EUR,
		AED,
		USDTinERC,
		USDTinTRC,
		USDCinERC,
	},
}

func GetExchangeCurrencies() map[Currency][]Currency {
	return ExchangedCurrencies
}

var CurrenciesByName = map[CurrencyName]Currency{
	// Fiat.
	RUBName: RUB,
	USDName: USD,
	EURName: EUR,
	AEDName: AED,
	UZSName: UZS,
	KZTName: KZT,
	CNYName: CNY,
	GBPName: GBP,
	TJSName: TJS,
	KGSName: KGS,
	GELName: GEL,
	MDLName: MDL,
	AMDName: AMD,

	// Crypto.
	USDTinERCName: USDTinERC,
	USDTinTRCName: USDTinTRC,
	USDCinERCName: USDCinERC,
	USDCinTRCName: USDCinTRC,
	BTCName:       BTC,
	ETHName:       ETH,
	TRXName:       TRX,
	TONName:       TON,
}

var CurrencyNamesByID = map[Currency]CurrencyName{
	// Fiat.
	RUB: RUBName,
	USD: USDName,
	EUR: EURName,
	AED: AEDName,
	UZS: UZSName,
	KZT: KZTName,
	CNY: CNYName,
	GBP: GBPName,
	TJS: TJSName,
	KGS: KGSName,
	GEL: GELName,
	MDL: MDLName,
	AMD: AMDName,

	// Crypto.
	USDTinERC: USDTinERCName,
	USDTinTRC: USDTinTRCName,
	USDCinERC: USDCinERCName,
	USDCinTRC: USDCinTRCName,
	BTC:       BTCName,
	ETH:       ETHName,
	TRX:       TRXName,
	TON:       TONName,
}
