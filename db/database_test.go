package db

import (
	"testing"
)

func TestFilterTravelInfoEligible_allFields(t *testing.T) {
	item1 := mockLocalDatabase()[0]
	item1.HasPassed = false
	item1.Place = "Warszawa"
	item1.Duration = ""
	item1.Distance = ""
	d := InitializeWith([]Competition{
		item1, anotherEligibleItemForTravel(),
	})
	items := d.FilterTravelInfoEligible()
	if len(items) != 2 {
		t.Error("Expected 2 item, got ", len(items))
	}
}

// Both Distance and Duration have to be empty to fetch travel info
func TestFilterTravelInfoEligible_noDistance(t *testing.T) {
	item1 := mockLocalDatabase()[0]
	item1.HasPassed = false
	item1.Place = "Warszawa"
	item1.Duration = "1h"
	item1.Distance = ""
	d := InitializeWith([]Competition{
		item1, anotherEligibleItemForTravel(),
	})
	items := d.FilterTravelInfoEligible()
	if len(items) != 1 {
		t.Error("Expected 1 item, got ", len(items))
	}
}

func TestFilterTravelInfoEligible_noDuration(t *testing.T) {
	item1 := mockLocalDatabase()[0]
	item1.HasPassed = false
	item1.Place = "Warszawa"
	item1.Duration = ""
	item1.Distance = "0km"
	d := InitializeWith([]Competition{
		item1, anotherEligibleItemForTravel(),
	})
	items := d.FilterTravelInfoEligible()
	if len(items) != 1 {
		t.Error("Expected 1 item, got ", len(items))
	}
}

func TestFilterTravelInfoEligible_onLineEvent(t *testing.T) {
	item1 := mockLocalDatabase()[0]
	item1.HasPassed = false
	item1.Place = "zawody online"
	item1.Duration = ""
	item1.Distance = ""
	d := InitializeWith([]Competition{
		item1, anotherEligibleItemForTravel(),
	})
	items := d.FilterTravelInfoEligible()
	if len(items) != 1 {
		t.Error("Expected 1 item, got ", len(items))
	}
}

func TestFilterTravelInfoEligible_passedEvent(t *testing.T) {
	item1 := mockLocalDatabase()[0]
	item1.HasPassed = true
	item1.Place = "Warszawa"
	item1.Duration = ""
	item1.Distance = ""
	d := InitializeWith([]Competition{
		item1, anotherEligibleItemForTravel(),
	})
	items := d.FilterTravelInfoEligible()
	if len(items) != 1 {
		t.Error("Expected 1 item, got ", len(items))
	}
}

func TestFilterWCAApiEligible_allFields(t *testing.T) {
	item1 := mockLocalDatabase()[0]
	item1.Type = CompetitionType.WCA
	item1.Events = []string{}
	item1.MainEvent = ""
	item1.CompetitorLimit = 0
	item1.Registered = 0
	d := InitializeWith([]Competition{
		item1, anotherEligibleItemForWCA(),
	})
	items := d.FilterWCAApiEligible()
	if len(items) != 2 {
		t.Error("Expected 2 item, got ", len(items))
	}
}

func TestFilterWCAApiEligible_wrongType(t *testing.T) {
	item1 := mockLocalDatabase()[0]
	item1.Type = CompetitionType.PPO
	item1.Events = []string{}
	item1.MainEvent = ""
	item1.CompetitorLimit = 0
	item1.Registered = 0
	d := InitializeWith([]Competition{
		item1, anotherEligibleItemForWCA(),
	})
	items := d.FilterWCAApiEligible()
	if len(items) != 1 {
		t.Error("Expected 1 item, got ", len(items))
	}
}

func TestFilterWCAApiEligible_alreadyHasEvents(t *testing.T) {
	item1 := mockLocalDatabase()[0]
	item1.Type = CompetitionType.WCA
	item1.Events = []string{"333"}
	item1.MainEvent = ""
	item1.CompetitorLimit = 0
	item1.Registered = 0
	d := InitializeWith([]Competition{
		item1, anotherEligibleItemForWCA(),
	})
	items := d.FilterWCAApiEligible()
	if len(items) != 2 {
		t.Error("Expected 2 item, got ", len(items))
	}
}

func TestFilterWCAApiEligible_alreadyHasMainEvent(t *testing.T) {
	item1 := mockLocalDatabase()[0]
	item1.Type = CompetitionType.WCA
	item1.Events = []string{}
	item1.MainEvent = "333"
	item1.CompetitorLimit = 0
	item1.Registered = 0
	d := InitializeWith([]Competition{
		item1, anotherEligibleItemForWCA(),
	})
	items := d.FilterWCAApiEligible()
	if len(items) != 2 {
		t.Error("Expected 2 item, got ", len(items))
	}
}

func TestFilterWCAApiEligible_alreadyHasCompetitorLimit(t *testing.T) {
	item1 := mockLocalDatabase()[0]
	item1.Type = CompetitionType.WCA
	item1.Events = []string{}
	item1.MainEvent = ""
	item1.CompetitorLimit = 1
	item1.Registered = 0
	d := InitializeWith([]Competition{
		item1, anotherEligibleItemForWCA(),
	})
	items := d.FilterWCAApiEligible()
	if len(items) != 2 {
		t.Error("Expected 2 item, got ", len(items))
	}
}

func TestFilterWCAApiEligible_alreadyHasRegistered(t *testing.T) {
	item1 := mockLocalDatabase()[0]
	item1.Type = CompetitionType.WCA
	item1.Events = []string{}
	item1.MainEvent = ""
	item1.CompetitorLimit = 0
	item1.Registered = 1
	d := InitializeWith([]Competition{
		item1, anotherEligibleItemForWCA(),
	})
	items := d.FilterWCAApiEligible()
	if len(items) != 2 {
		t.Error("Expected 2 item, got ", len(items))
	}
}

func TestFilterScrapEligible_Cube4Fun(t *testing.T) {
	item1 := mockLocalDatabase()[0]
	item2 := mockLocalDatabase()[1]
	item3 := mockLocalDatabase()[3]
	item1.Type = CompetitionType.PPO
	item1.URL = "http://google.com"
	item2.Type = CompetitionType.Cube4Fun
	item2.URL = "http://google.com"
	item3.Type = CompetitionType.WCA
	item3.URL = "http://google.com"
	d := InitializeWith([]Competition{
		item1, item2,
	})
	items := d.FilterScrapCube4FunEligible()
	if len(items) != 1 {
		t.Error("Expected 1 items, got ", len(items))
	}
}

func TestFilterScrapEligible_Cube4FunWithoutURL(t *testing.T) {
	item1 := mockLocalDatabase()[0]
	item2 := mockLocalDatabase()[1]
	item3 := mockLocalDatabase()[3]
	item1.Type = CompetitionType.PPO
	item1.URL = "http://google.com"
	item2.Type = CompetitionType.Cube4Fun
	item2.URL = ""
	item3.Type = CompetitionType.WCA
	item3.URL = "http://google.com"
	d := InitializeWith([]Competition{
		item1, item2,
	})
	items := d.FilterScrapCube4FunEligible()
	if len(items) != 0 {
		t.Error("Expected 0 items, got ", len(items))
	}
}

func TestFilterScrapEligible_PPO(t *testing.T) {
	item1 := mockLocalDatabase()[0]
	item2 := mockLocalDatabase()[1]
	item3 := mockLocalDatabase()[3]
	item1.Type = CompetitionType.PPO
	item1.URL = "http://google.com"
	item2.Type = CompetitionType.Cube4Fun
	item2.URL = "http://google.com"
	item3.Type = CompetitionType.WCA
	item3.URL = "http://google.com"
	d := InitializeWith([]Competition{
		item1, item2,
	})
	items := d.FilterScrapPPOEligible()
	if len(items) != 1 {
		t.Error("Expected 1 items, got ", len(items))
	}
}

func TestFilterScrapEligible_PPOWithoutURL(t *testing.T) {
	item1 := mockLocalDatabase()[0]
	item2 := mockLocalDatabase()[1]
	item3 := mockLocalDatabase()[3]
	item1.Type = CompetitionType.PPO
	item1.URL = ""
	item2.Type = CompetitionType.Cube4Fun
	item2.URL = "http://google.com"
	item3.Type = CompetitionType.WCA
	item3.URL = "http://google.com"
	d := InitializeWith([]Competition{
		item1, item2,
	})
	items := d.FilterScrapPPOEligible()
	if len(items) != 0 {
		t.Error("Expected 0 items, got ", len(items))
	}
}

func anotherEligibleItemForTravel() Competition {
	item := mockLocalDatabase()[1]
	item.HasPassed = false
	item.Place = "Warszawa"
	item.Duration = ""
	item.Distance = ""

	return item
}

func anotherEligibleItemForWCA() Competition {
	item := mockLocalDatabase()[1]
	item.Type = CompetitionType.WCA
	item.Events = []string{}
	item.MainEvent = ""
	item.CompetitorLimit = 0
	item.Registered = 0

	return item
}
