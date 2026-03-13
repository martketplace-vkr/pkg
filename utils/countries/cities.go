package countries

type CityID int

const (
	CityWalletID = 0

	CityRUMoscowID          CityID = 1
	CityRUSaintPetersburgID CityID = 2
	CityRUNovosibirskID     CityID = 3
	CityRUKrasnodarID       CityID = 4
	CityRUKrasnoyarskID     CityID = 5
	CityRUSochiID           CityID = 6
	CityRUKazanID           CityID = 7
	CityRUEkaterinburgID    CityID = 8
	CityRURostovID          CityID = 9
	CityRUStavropolID       CityID = 10
	CityRUUfaID             CityID = 11
	CityRUNizhnyNovgorodID  CityID = 12
	CityRUVoronezhID        CityID = 13
	CityRUKaliningradID     CityID = 14

	CityRUIzhevskID          CityID = 15
	CityRUSamaraID           CityID = 16
	CityRUChelyabinskID      CityID = 17
	CityRUVladivostokID      CityID = 18
	CityRUYuzhnoSakhalinskID CityID = 19
	CityRUVolgogradID        CityID = 20
	CityRUOmskID             CityID = 21
	CityRUHabarovskID        CityID = 22
	CityRUBarbaulID          CityID = 23
	CityRUPenzaID            CityID = 24
	CityRUPermID             CityID = 25
	CityRUSaratovID          CityID = 26
	CityRUTyumenID           CityID = 27
	CityRULipetskID          CityID = 28
	CityRUMakhachkalaID      CityID = 29
	CityRURyazanID           CityID = 30
	CityRUKemerovoID         CityID = 31
	CityRUBiyskID            CityID = 32
	CityRUGornoAltayskID     CityID = 33
	CityRUAstrakhanID        CityID = 34
	CityRUUlyanovskID        CityID = 35
	CityRUIrkutskID          CityID = 36
	CityRUYaroslavlID        CityID = 37
	CityRUTomskID            CityID = 38

	CityAEDubaiID    CityID = 1
	CityAEAbuDhabiID CityID = 2
	CityKZAstanaID   CityID = 1
	CityKZAlmatyID   CityID = 2
	CityBYMinskID    CityID = 1
	CityUAKievID     CityID = 1
	CityGETbilisiID  CityID = 1
	CityTRIstanbulID CityID = 1
	CityTRAntalyaID  CityID = 2
	CityAMYerevanID  CityID = 1
	CityKGBishkekID  CityID = 1
)

type City struct {
	ID    CityID
	Names map[Locale]string
}

func (c City) Locale(loc Locale) string {
	if locale, ok := c.Names[loc]; !ok {
		return ""
	} else {
		return locale
	}
}

var (
	CityRUMoscow = City{
		ID: CityRUMoscowID,
		Names: map[Locale]string{
			RU: "Москва",
			EN: "Moscow",
		},
	}
	CityRUSaintPetersburg = City{
		ID: CityRUSaintPetersburgID,
		Names: map[Locale]string{
			RU: "Санкт-Петербург",
			EN: "Saint Petersburg",
		},
	}
	CityRUNovosibirsk = City{
		ID: CityRUNovosibirskID,
		Names: map[Locale]string{
			RU: "Новосибирск",
			EN: "Novosibirsk",
		},
	}
	CityRUKrasnodar = City{
		ID: CityRUKrasnodarID,
		Names: map[Locale]string{
			RU: "Краснодар",
			EN: "Krasnodar",
		},
	}
	CityRUKrasnoyarsk = City{
		ID: CityRUKrasnoyarskID,
		Names: map[Locale]string{
			RU: "Красноярск",
			EN: "Krasnoyarsk",
		},
	}
	CityRUSochi = City{
		ID: CityRUSochiID,
		Names: map[Locale]string{
			RU: "Сочи",
			EN: "Sochi",
		},
	}
	CityRUKazan = City{
		ID: CityRUKazanID,
		Names: map[Locale]string{
			RU: "Казань",
			EN: "Kazan",
		},
	}
	CityRUEkaterinburg = City{
		ID: CityRUEkaterinburgID,
		Names: map[Locale]string{
			RU: "Екатеринбург",
			EN: "Ekaterinburg",
		},
	}
	CityRURostov = City{
		ID: CityRURostovID,
		Names: map[Locale]string{
			RU: "Ростов-на-Дону",
			EN: "Rostov-on-Don",
		},
	}
	CityRUStavropol = City{
		ID: CityRUStavropolID,
		Names: map[Locale]string{
			RU: "Ставрополь",
			EN: "Stavropol",
		},
	}
	CityRUUfa = City{
		ID: CityRUUfaID,
		Names: map[Locale]string{
			RU: "Уфа",
			EN: "Ufa",
		},
	}
	CityRUNizhnyNovgorod = City{
		ID: CityRUNizhnyNovgorodID,
		Names: map[Locale]string{
			RU: "Нижний Новгород",
			EN: "Nizhny Novgorod",
		},
	}
	CityRUVoronezh = City{
		ID: CityRUVoronezhID,
		Names: map[Locale]string{
			RU: "Воронеж",
			EN: "Voronezh",
		},
	}
	CityRUKaliningrad = City{
		ID: CityRUKaliningradID,
		Names: map[Locale]string{
			RU: "Калининград",
			EN: "Kaliningrad",
		},
	}
	CityRUIzhevsk = City{
		ID: CityRUIzhevskID,
		Names: map[Locale]string{
			RU: "Ижевск",
			EN: "Izhevsk",
		},
	}
	CityRUSamara = City{
		ID: CityRUSamaraID,
		Names: map[Locale]string{
			RU: "Самара",
			EN: "Samara",
		},
	}
	CityRUChelyabinsk = City{
		ID: CityRUChelyabinskID,
		Names: map[Locale]string{
			RU: "Челябинск",
			EN: "Chelyabinsk",
		},
	}
	CityRUVladivostok = City{
		ID: CityRUVladivostokID,
		Names: map[Locale]string{
			RU: "Владивосток",
			EN: "Vladivostok",
		},
	}
	CityRUYuzhnoSakhalinsk = City{
		ID: CityRUYuzhnoSakhalinskID,
		Names: map[Locale]string{
			RU: "Южно-Сахалинск",
			EN: "Yuzhno-Sakhalinsk",
		},
	}
	CityRUVolgograd = City{
		ID: CityRUVolgogradID,
		Names: map[Locale]string{
			RU: "Волгоград",
			EN: "Volgograd",
		},
	}
	CityRUOmsk = City{
		ID: CityRUOmskID,
		Names: map[Locale]string{
			RU: "Омск",
			EN: "Omsk",
		},
	}
	CityRUHabarovsk = City{
		ID: CityRUHabarovskID,
		Names: map[Locale]string{
			RU: "Хабаровск",
			EN: "Habarovsk",
		},
	}
	CityRUBarbaul = City{
		ID: CityRUBarbaulID,
		Names: map[Locale]string{
			RU: "Барнаул",
			EN: "Barbaul",
		},
	}
	CityRUPenza = City{
		ID: CityRUPenzaID,
		Names: map[Locale]string{
			RU: "Пенза",
			EN: "Penza",
		},
	}
	CityRUPerm = City{
		ID: CityRUPermID,
		Names: map[Locale]string{
			RU: "Пермь",
			EN: "Perm",
		},
	}
	CityRUSaratov = City{
		ID: CityRUSaratovID,
		Names: map[Locale]string{
			RU: "Саратов",
			EN: "Saratov",
		},
	}
	CityRUTyumen = City{
		ID: CityRUTyumenID,
		Names: map[Locale]string{
			RU: "Тюмень",
			EN: "Tyumen",
		},
	}
	CityRULipetsk = City{
		ID: CityRULipetskID,
		Names: map[Locale]string{
			RU: "Липецк",
			EN: "Lipetsk",
		},
	}
	CityRUMakhachkala = City{
		ID: CityRUMakhachkalaID,
		Names: map[Locale]string{
			RU: "Махачкала",
			EN: "Makhachkala",
		},
	}
	CityRURyazan = City{
		ID: CityRURyazanID,
		Names: map[Locale]string{
			RU: "Рязань",
			EN: "Ryazan",
		},
	}
	CityRUKemerovo = City{
		ID: CityRUKemerovoID,
		Names: map[Locale]string{
			RU: "Кемерово",
			EN: "Kemerovo",
		},
	}
	CityRUBiysk = City{
		ID: CityRUBiyskID,
		Names: map[Locale]string{
			RU: "Бийск",
			EN: "Biysk",
		},
	}
	CityRUGornoAltaysk = City{
		ID: CityRUGornoAltayskID,
		Names: map[Locale]string{
			RU: "Горно-Алтайск",
			EN: "Gorno-Altaysk",
		},
	}
	CityRUAstrakhan = City{
		ID: CityRUAstrakhanID,
		Names: map[Locale]string{
			RU: "Астрахань",
			EN: "Astrakhan",
		},
	}
	CityRUUlyanovsk = City{
		ID: CityRUUlyanovskID,
		Names: map[Locale]string{
			RU: "Ульяновск",
			EN: "Ulyanovsk",
		},
	}
	CityRUIrkutsk = City{
		ID: CityRUIrkutskID,
		Names: map[Locale]string{
			RU: "Иркутск",
			EN: "Irkutsk",
		},
	}
	CityRUYaroslavl = City{
		ID: CityRUYaroslavlID,
		Names: map[Locale]string{
			RU: "Ярославль",
			EN: "Yaroslavl",
		},
	}
	CityRUTomsk = City{
		ID: CityRUTomskID,
		Names: map[Locale]string{
			RU: "Томск",
			EN: "Tomsk",
		},
	}

	CityAEDubai = City{
		ID: CityAEDubaiID,
		Names: map[Locale]string{
			RU: "Дубай",
			EN: "Dubai",
		},
	}

	CityAEAbuDhabi = City{
		ID: CityAEAbuDhabiID,
		Names: map[Locale]string{
			RU: "Абу-Даби",
			EN: "Abu Dhabi",
		},
	}

	CityKZAstana = City{
		ID: CityKZAstanaID,
		Names: map[Locale]string{
			RU: "Астана",
			EN: "Astana",
		},
	}

	CityKZAlmaty = City{
		ID: CityKZAlmatyID,
		Names: map[Locale]string{
			RU: "Алматы",
			EN: "Almaty",
		},
	}

	CityBYMinsk = City{
		ID: CityBYMinskID,
		Names: map[Locale]string{
			RU: "Минск",
			EN: "Minsk",
		},
	}

	CityUAKiev = City{
		ID: CityUAKievID,
		Names: map[Locale]string{
			RU: "Киев",
			EN: "Kiev",
		},
	}

	CityGETbilisi = City{
		ID: CityGETbilisiID,
		Names: map[Locale]string{
			RU: "Тбилиси",
			EN: "Tbilisi",
		},
	}

	CityTRIstanbul = City{
		ID: CityTRIstanbulID,
		Names: map[Locale]string{
			RU: "Стамбул",
			EN: "Istanbul",
		},
	}

	CityTRAntalya = City{
		ID: CityTRAntalyaID,
		Names: map[Locale]string{
			RU: "Анталья",
			EN: "Antalya",
		},
	}

	CityAMYerevan = City{
		ID: CityAMYerevanID,
		Names: map[Locale]string{
			RU: "Ереван",
			EN: "Yerevan",
		},
	}

	CityKGBishkek = City{
		ID: CityKGBishkekID,
		Names: map[Locale]string{
			RU: "Бишкек",
			EN: "Bishkek",
		},
	}
)
