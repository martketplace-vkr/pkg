package countries

import (
	"errors"

	"github.com/martketplace-vkr/pkg/utils/currency"
)

var (
	ErrInvalidCountry = errors.New("invalid country")
	ErrInvalidCity    = errors.New("invalid city")
	ErrInvalidLocale  = errors.New("invalid locale")
)

type (
	Location struct {
		CountryID CountryID
		CityID    CityID
	}

	Coupling struct {
		Default   Location
		Available []Location
	}
)

var CurrenciesAvailability = map[currency.Currency]Coupling{
	currency.USD: {
		Default: Location{
			CountryID: CountryRUID,
			CityID:    CityRUMoscowID,
		}, Available: []Location{
			{CountryID: CountryRUID, CityID: CityRUMoscowID},
			{CountryID: CountryRUID, CityID: CityRUSaintPetersburgID},
			{CountryID: CountryRUID, CityID: CityRUNovosibirskID},
			{CountryID: CountryRUID, CityID: CityRUKrasnodarID},
			{CountryID: CountryRUID, CityID: CityRUKrasnoyarskID},
			{CountryID: CountryRUID, CityID: CityRUSochiID},
			{CountryID: CountryRUID, CityID: CityRUKazanID},
			{CountryID: CountryRUID, CityID: CityRUEkaterinburgID},
			{CountryID: CountryRUID, CityID: CityRURostovID},
			{CountryID: CountryAMID, CityID: CityAMYerevanID},
			{CountryID: CountryKGID, CityID: CityKGBishkekID},
			{CountryID: CountryGEID, CityID: CityGETbilisiID},
			{CountryID: CountryAEID, CityID: CityAEDubaiID},
			{CountryID: CountryRUID, CityID: CityRUKaliningradID},
			{CountryID: CountryRUID, CityID: CityRUNizhnyNovgorodID},
			{CountryID: CountryRUID, CityID: CityRUVoronezhID},
			{CountryID: CountryRUID, CityID: CityRUStavropolID},
			{CountryID: CountryRUID, CityID: CityRUUfaID},
			{CountryID: CountryRUID, CityID: CityRUHabarovskID},
			{CountryID: CountryRUID, CityID: CityRUBarbaulID},
			{CountryID: CountryRUID, CityID: CityRUPenzaID},
			{CountryID: CountryRUID, CityID: CityRUPermID},
			{CountryID: CountryRUID, CityID: CityRUSaratovID},
		}},
	currency.RUB: {
		Default: Location{
			CountryID: CountryRUID,
			CityID:    CityRUMoscowID,
		}, Available: []Location{
			{CountryID: CountryRUID, CityID: CityRUMoscowID},
			{CountryID: CountryRUID, CityID: CityRUKrasnodarID},
			{CountryID: CountryRUID, CityID: CityRUNovosibirskID},
			{CountryID: CountryRUID, CityID: CityRUSaintPetersburgID},
			{CountryID: CountryRUID, CityID: CityRUKrasnoyarskID},
			{CountryID: CountryRUID, CityID: CityRUSochiID},
			{CountryID: CountryRUID, CityID: CityRUKazanID},
			{CountryID: CountryRUID, CityID: CityRUEkaterinburgID},
			{CountryID: CountryRUID, CityID: CityRURostovID},
			{CountryID: CountryRUID, CityID: CityRUKaliningradID},
			{CountryID: CountryRUID, CityID: CityRUNizhnyNovgorodID},
			{CountryID: CountryRUID, CityID: CityRUVoronezhID},
			{CountryID: CountryRUID, CityID: CityRUStavropolID},
			{CountryID: CountryRUID, CityID: CityRUUfaID},
			{CountryID: CountryRUID, CityID: CityRUIzhevskID},
			{CountryID: CountryRUID, CityID: CityRUSamaraID},
			{CountryID: CountryRUID, CityID: CityRUChelyabinskID},
			{CountryID: CountryRUID, CityID: CityRUVladivostokID},
			{CountryID: CountryRUID, CityID: CityRUYuzhnoSakhalinskID},
			{CountryID: CountryRUID, CityID: CityRUVolgogradID},
			{CountryID: CountryRUID, CityID: CityRUOmskID},
			{CountryID: CountryRUID, CityID: CityRUHabarovskID},
			{CountryID: CountryRUID, CityID: CityRUBarbaulID},
			{CountryID: CountryRUID, CityID: CityRUPenzaID},
			{CountryID: CountryRUID, CityID: CityRUPermID},
			{CountryID: CountryRUID, CityID: CityRUSaratovID},
			{CountryID: CountryRUID, CityID: CityRUTyumenID},
			{CountryID: CountryRUID, CityID: CityRULipetskID},
			{CountryID: CountryRUID, CityID: CityRUMakhachkalaID},
			{CountryID: CountryRUID, CityID: CityRURyazanID},
			{CountryID: CountryRUID, CityID: CityRUKemerovoID},
			{CountryID: CountryRUID, CityID: CityRUBiyskID},
			{CountryID: CountryRUID, CityID: CityRUGornoAltayskID},
			{CountryID: CountryRUID, CityID: CityRUAstrakhanID},
			{CountryID: CountryRUID, CityID: CityRUUlyanovskID},
			{CountryID: CountryRUID, CityID: CityRUIrkutskID},
			{CountryID: CountryKGID, CityID: CityKGBishkekID},
			{CountryID: CountryRUID, CityID: CityRUYaroslavlID},
			{CountryID: CountryRUID, CityID: CityRUTomskID},
		}},
	currency.TRX: {
		Default: Location{
			CountryID: CountryRUID,
			CityID:    CityRUMoscowID,
		}, Available: []Location{
			{CountryID: CountryRUID, CityID: CityRUMoscowID},
		},
	},
	currency.ETH: {
		Default: Location{
			CountryID: CountryRUID,
			CityID:    CityRUMoscowID,
		}, Available: []Location{
			{CountryID: CountryRUID, CityID: CityRUMoscowID},
		},
	},
	currency.BTC: {
		Default: Location{
			CountryID: CountryRUID,
			CityID:    CityRUMoscowID,
		}, Available: []Location{
			{CountryID: CountryRUID, CityID: CityRUMoscowID},
		},
	},
	currency.USDTinTRC: {
		Default: Location{
			CountryID: CountryRUID,
			CityID:    CityRUMoscowID,
		}, Available: []Location{
			{CountryID: CountryRUID, CityID: CityRUMoscowID},
		},
	},
	currency.USDTinERC: {
		Default: Location{
			CountryID: CountryRUID,
			CityID:    CityRUMoscowID,
		}, Available: []Location{
			{CountryID: CountryRUID, CityID: CityRUMoscowID},
		},
	},
	currency.AED: {
		Default: Location{
			CountryID: CountryAEID,
			CityID:    CityAEDubaiID,
		}, Available: []Location{
			{CountryID: CountryAEID, CityID: CityAEDubaiID},
		},
	},
	currency.EUR: {
		Default: Location{
			CountryID: CountryRUID,
			CityID:    CityRUMoscowID,
		}, Available: []Location{
			{CountryID: CountryRUID, CityID: CityRUMoscowID},
			{CountryID: CountryRUID, CityID: CityRUSaintPetersburgID},
		},
	},
}
