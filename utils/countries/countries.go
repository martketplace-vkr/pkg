package countries

type CountryID int

// ! BEWARE: following constants shall never be > 999.
const (
	CountryWallet CountryID = 0
	CountryRUID   CountryID = 1
	CountryAEID   CountryID = 2
	CountryKZID   CountryID = 3
	CountryBYID   CountryID = 4
	CountryUAID   CountryID = 5
	CountryGEID   CountryID = 6
	CountryTRID   CountryID = 7
	CountryAMID   CountryID = 8
	CountryKGID   CountryID = 9
)

var Countries = map[CountryID][]City{
	CountryRUID: CountryRU.Cities,
	CountryAEID: CountryAE.Cities,
	CountryKZID: CountryKZ.Cities,
	CountryGEID: CountryGE.Cities,
	CountryTRID: CountryTR.Cities,
	CountryAMID: CountryAM.Cities,
	CountryKGID: CountryKG.Cities,
}

type Country struct {
	ID     CountryID
	Names  map[Locale]string
	Cities []City
}

func (c Country) Locale(loc Locale) string {
	if locale, ok := c.Names[loc]; !ok {
		return ""
	} else {
		return locale
	}
}

var (
	CountryRU = Country{
		ID: CountryRUID,
		Names: map[Locale]string{
			RU: "Россия",
			EN: "Russia",
		},
		Cities: []City{
			CityRUMoscow,
			CityRUSaintPetersburg,
			CityRUNovosibirsk,
			CityRUKrasnodar,
			CityRUKrasnoyarsk,
			CityRUSochi,
			CityRUKazan,
			CityRUEkaterinburg,
			CityRURostov,
			CityRUStavropol,
			CityRUUfa,
			CityRUNizhnyNovgorod,
			CityRUVoronezh,
			CityRUKaliningrad,
			CityRUIzhevsk,
			CityRUSamara,
			CityRUChelyabinsk,
			CityRUVladivostok,
			CityRUYuzhnoSakhalinsk,
			CityRUVolgograd,
			CityRUOmsk,
			CityRUHabarovsk,
			CityRUBarbaul,
			CityRUPenza,
			CityRUPerm,
			CityRUSaratov,
			CityRUTyumen,
			CityRULipetsk,
			CityRUMakhachkala,
			CityRURyazan,
			CityRUKemerovo,
			CityRUBiysk,
			CityRUGornoAltaysk,
			CityRUAstrakhan,
			CityRUUlyanovsk,
			CityRUIrkutsk,
			CityRUYaroslavl,
			CityRUTomsk,
		},
	}
	CountryAE = Country{
		ID: CountryAEID,
		Names: map[Locale]string{
			RU: "ОАЭ",
			EN: "UAE",
		},
		Cities: []City{
			CityAEDubai,
			CityAEAbuDhabi,
		},
	}
	CountryKZ = Country{
		ID: CountryKZID,
		Names: map[Locale]string{
			RU: "Казахстан",
			EN: "Kazakhstan",
		},
		Cities: []City{
			CityKZAstana,
			CityKZAlmaty,
		},
	}
	CountryGE = Country{
		ID: CountryGEID,
		Names: map[Locale]string{
			RU: "Грузия",
			EN: "Georgia",
		},
		Cities: []City{
			CityGETbilisi,
		},
	}
	CountryTR = Country{
		ID: CountryTRID,
		Names: map[Locale]string{
			RU: "Турция",
			EN: "Turkey",
		},
		Cities: []City{
			CityTRIstanbul,
			CityTRAntalya,
		},
	}
	CountryAM = Country{
		ID: CountryAMID,
		Names: map[Locale]string{
			RU: "Армения",
			EN: "Armenia",
		},
		Cities: []City{
			CityAMYerevan,
		},
	}
	CountryKG = Country{
		ID: CountryKGID,
		Names: map[Locale]string{
			RU: "Киргизия",
			EN: "Kyrgyzstan",
		},
		Cities: []City{
			CityKGBishkek,
		},
	}
	CountryBY = Country{
		ID: CountryBYID,
		Names: map[Locale]string{
			RU: "Беларусь",
			EN: "Belarus",
		},
		Cities: []City{
			CityBYMinsk,
		},
	}
	CountryUA = Country{
		ID: CountryUAID,
		Names: map[Locale]string{
			RU: "Украина",
			EN: "Ukraine",
		},
		Cities: []City{
			CityUAKiev,
		},
	}
)
