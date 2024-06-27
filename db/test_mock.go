package db

import (
	"math/rand"
)

func randomItem(items []string) string {
	if len(items) == 0 {
		return ""
	}
	randomIndex := rand.Intn(len(items))
	return items[randomIndex]
}

func randomEvents() []string {
	items := []string{
		"sq1",
		"skewb",
		"pyram",
		"minx",
		"222",
		"333bf",
		"333",
		"333fm",
		"333mbf",
		"333mbo",
		"333oh",
		"333ft",
		"444bf",
		"444",
		"555bf",
		"555",
		"666",
		"777",
		"clock",
		"magic",
	}
	if len(items) == 0 {
		return []string{""}
	}

	subsetLength := rand.Intn(len(items)) + 1

	rand.Shuffle(len(items), func(i, j int) {
		items[i], items[j] = items[j], items[i]
	})

	return items[:subsetLength]
}

func mockLocalDatabase() []Competition {
	return []Competition{
		{
			Type:            "Unknown",
			TypeSpecificId:  "",
			WCAId:           "",
			Id:              "mistrzostwa-polski-w-speedcubingu-2024",
			Header:          "Listopad 2024",
			Name:            "[Has URL] Mistrzostwa Polski w Speedcubingu 2024",
			URL:             "https://zawody.rubiart.pl/competitions/48",
			Place:           "Warszawa",
			LogoURL:         "https://www.speedcubing.pl/content/small/logo-pn-2023.png",
			ContactName:     "Zespół organizacyjny MP24",
			ContactURL:      "foobar@liganiezwyklych.pl",
			Date:            "8-10 listopada 2024 r.",
			Distance:        "",
			Duration:        "",
			HasWCA:          true,
			HasPassed:       false,
			Events:          []string{},
			MainEvent:       "",
			CompetitorLimit: 0,
			Registered:      0,
		},
		{
			Type:            "Unknown",
			TypeSpecificId:  "",
			WCAId:           "",
			Id:              "brizzon-sylwester-open-2024",
			Header:          "Grudzień 2024",
			Name:            "[No URL] Brizzon Sylwester Open 2024",
			URL:             "",
			Place:           "Poznań",
			LogoURL:         "https://www.speedcubing.pl/content/small/koziolki-logo.jpg",
			ContactName:     "Zespół Organizacyjny",
			ContactURL:      "foobar@googlegroups.com",
			Date:            "31 grudnia 2024 r.",
			Distance:        "",
			Duration:        "",
			HasWCA:          true,
			HasPassed:       false,
			Events:          []string{},
			MainEvent:       "",
			CompetitorLimit: 0,
			Registered:      0,
		},
		{
			Type:            "Cube4Fun",
			TypeSpecificId:  "bydgoszcz-2024",
			WCAId:           "",
			Id:              "cube4fun-biocube-bydgoszcz-2024",
			Header:          "Lipiec 2024",
			Name:            "[Has URL] Cube4fun Biocube Bydgoszcz 2024",
			URL:             "https://cube4fun.pl/competition/bydgoszcz-2024",
			Place:           "Bydgoszcz",
			LogoURL:         "https://www.speedcubing.pl/content/small/bydgoszcz-c4f.png",
			ContactName:     "cube4fun",
			ContactURL:      "foobar@cube4fun.pl",
			Date:            "6-7 lipca 2024 r.",
			Distance:        "",
			Duration:        "",
			HasWCA:          true,
			HasPassed:       false,
			Events:          []string{},
			MainEvent:       "",
			CompetitorLimit: 0,
			Registered:      0,
		},
		{
			Type:            "Cube4Fun",
			TypeSpecificId:  "",
			WCAId:           "",
			Id:              "cube4fun-in-biala-podlaska-2024",
			Header:          "Lipiec 2024",
			Name:            "[No URL] Cube4fun in Biała Podlaska 2024",
			URL:             "",
			Place:           "Biała Podlaska",
			LogoURL:         "https://www.speedcubing.pl/content/small/biala-podlasaka-2024.png",
			ContactName:     "Cube4fun",
			ContactURL:      "foobar@cube4fun.pl",
			Date:            "12-14 lipca 2024 r.",
			Distance:        "",
			Duration:        "",
			HasWCA:          true,
			HasPassed:       false,
			Events:          []string{},
			MainEvent:       "",
			CompetitorLimit: 0,
			Registered:      0,
		},
		{
			Type:            "WCA",
			TypeSpecificId:  "ZoryCubingMansion2024",
			WCAId:           "",
			Id:              "zory-cubing-mansion-2024",
			Header:          "Lipiec 2024",
			Name:            "[Has URL] Żory Cubing Mansion 2024",
			URL:             "https://www.worldcubeassociation.org/competitions/ZoryCubingMansion2024",
			Place:           "Żory",
			LogoURL:         "https://www.speedcubing.pl/content/small/zory-cubing-mansion.png",
			ContactName:     "Karol Zakrzewski",
			ContactURL:      "foobar@worldcubeassociation.org",
			Date:            "27 lipca 2024 r.",
			Distance:        "",
			Duration:        "",
			HasWCA:          true,
			HasPassed:       false,
			Events:          []string{},
			MainEvent:       "",
			CompetitorLimit: 0,
			Registered:      0,
		},
		{
			Type:            "Unknown",
			TypeSpecificId:  "",
			WCAId:           "",
			Id:              "ppo-iv-2024",
			Header:          "Wrzesień 2024",
			Name:            "[No URL] PPO IV 2024",
			URL:             "",
			Place:           "zawody online",
			LogoURL:         "https://www.speedcubing.pl/content/small/ppo-2024.png",
			ContactName:     "Adam Polkowski",
			ContactURL:      "adam.polkowski@rubiart.pl",
			Date:            "2-11 września 2024 r.",
			Distance:        "",
			Duration:        "",
			HasWCA:          false,
			HasPassed:       false,
			Events:          []string{},
			MainEvent:       "",
			CompetitorLimit: 0,
			Registered:      0,
		},
		{
			Type:            "PPO",
			TypeSpecificId:  "49",
			WCAId:           "",
			Id:              "ppo-final-2024",
			Header:          "Listopad 2024",
			Name:            "[Has URL] PPO Final 2024",
			URL:             "https://zawody.rubiart.pl/competitions/49",
			Place:           "zawody online",
			LogoURL:         "https://www.speedcubing.pl/content/small/ppo-2024.png",
			ContactName:     "Adam Polkowski",
			ContactURL:      "adam.polkowski@rubiart.pl",
			Date:            "18-27 listopada 2024 r.",
			Distance:        "",
			Duration:        "",
			HasWCA:          false,
			HasPassed:       false,
			Events:          []string{},
			MainEvent:       "",
			CompetitorLimit: 0,
			Registered:      0,
		},
		{
			Type:            "Unknown",
			TypeSpecificId:  "",
			WCAId:           "",
			Id:              "gls-v-2024",
			Header:          "Październik 2024",
			Name:            "[No URL] GLS V 2024",
			URL:             "",
			Place:           "Gdańsk",
			LogoURL:         "https://www.speedcubing.pl/content/small/logo-gls-2023.png",
			ContactName:     "Adam Polkowski",
			ContactURL:      "adam.polkowski@rubiart.pl",
			Date:            "5-6 października 2024 r.",
			Distance:        "",
			Duration:        "",
			HasWCA:          true,
			HasPassed:       false,
			Events:          []string{},
			MainEvent:       "",
			CompetitorLimit: 0,
			Registered:      0,
		},
		{
			Type:            "Unknown",
			TypeSpecificId:  "",
			WCAId:           "",
			Id:              "mistrzostwa-polski-juniorow-2024",
			Header:          "Październik 2024",
			Name:            "Mistrzostwa Polski Juniorów 2024",
			URL:             "",
			Place:           "zawody online",
			LogoURL:         "https://www.speedcubing.pl/content/small/logo-mpj2020.png",
			ContactName:     "Adam Polkowski",
			ContactURL:      "adam.polkowski@rubiart.pl",
			Date:            "14-24 października 2024 r.",
			Distance:        "",
			Duration:        "",
			HasWCA:          false,
			HasPassed:       false,
			Events:          []string{},
			MainEvent:       "",
			CompetitorLimit: 0,
			Registered:      0,
		},
	}
}

var possibleType = []string{
	CompetitionType.WCA,
	CompetitionType.Unknown,
	CompetitionType.PPO,
	CompetitionType.Cube4Fun,
}
var possibleTypeSpecificId = []string{
	"",
}
var possibleWCAId = []string{}
var possibleId = []string{}
var possibleHeader = []string{}
var possibleName = []string{}
var possibleURL = []string{}
var possiblePlace = []string{}
var possibleLogoURL = []string{}
var possibleContactName = []string{}
var possibleContactURL = []string{}
var possibleDate = []string{}
var possibleDistance = []string{}
var possibleDuration = []string{}
var possibleHasWCA = []string{}
var possibleHasPassed = []string{}
var possibleEvents = []string{}
var possibleMainEvent = []string{}
var possibleCompetitorLimit = []string{}
var possibleRegistered = []string{}
