package db

import (
	"testing"
)

func TestCreation(t *testing.T) {
	d := InitializeWith([]Competition{})
	if len(d.GetAll()) != 0 {
		t.Error("Expected 0 items, got ", len(d.GetAll()))
	}
}

func TestAddItems(t *testing.T) {
	l := InitializeWith([]Competition{mockLocalDatabase()[0]})
	d := InitializeWith([]Competition{
		mockLocalDatabase()[1],
		mockLocalDatabase()[2],
	})
	m := Merger{
		added:   []string{"mistrzostwa-polski-w-speedcubingu-2024"},
		passed:  []string{},
		changed: []string{},
	}
	n := m.Merge(l, d)
	if len(n.GetAll()) != 3 {
		t.Error("Expected 3 items, got ", len(d.GetAll()))
	}
}

func TestMarkAsPassedItems(t *testing.T) {
	l := InitializeWith([]Competition{})
	d := InitializeWith([]Competition{
		mockLocalDatabase()[1],
		mockLocalDatabase()[2],
	})
	m := Merger{
		added:   []string{},
		passed:  []string{"brizzon-sylwester-open-2024"},
		changed: []string{},
	}
	n := m.Merge(l, d)
	item := n.Get("brizzon-sylwester-open-2024")
	if !item.HasPassed {
		t.Error("Expected hasPassed = true, got ", item.HasPassed)
	}
}
func TestChangeItems(t *testing.T) {
	item1 := mockLocalDatabase()[1]
	item2 := mockLocalDatabase()[2]
	localItem1 := item1
	localItem1.URL = "New URL"
	l := InitializeWith([]Competition{localItem1})
	d := InitializeWith([]Competition{
		item1, item2,
	})
	m := Merger{
		added:   []string{},
		passed:  []string{},
		changed: []string{"brizzon-sylwester-open-2024"},
	}
	n := m.Merge(l, d)
	item := n.Get("brizzon-sylwester-open-2024")
	if item.URL != "New URL" {
		t.Error("Expected change in the URL, got ", item.URL)
	}
}

func TestAllChangesAtOnceItems(t *testing.T) {
	item1 := mockLocalDatabase()[1]
	item2 := mockLocalDatabase()[2]
	item3 := mockLocalDatabase()[3]
	localItem1 := item1
	localItem1.URL = "New URL"
	localItem2 := item2
	localItem2.URL = "New URL"
	l := InitializeWith([]Competition{
		item3, localItem1, localItem2,
	})
	d := InitializeWith([]Competition{
		item1, item2,
	})
	m := Merger{
		added:   []string{"cube4fun-in-biala-podlaska-2024"},
		passed:  []string{"cube4fun-biocube-bydgoszcz-2024"},
		changed: []string{"brizzon-sylwester-open-2024", "cube4fun-biocube-bydgoszcz-2024"},
	}
	n := m.Merge(l, d)

	if len(n.GetAll()) != 3 {
		t.Error("Expected 3 items, got ", len(d.GetAll()))
	}

	item := n.Get("brizzon-sylwester-open-2024")
	if item.URL != "New URL" {
		t.Error("Expected change in the URL, got ", item.URL)
	}

	item = n.Get("cube4fun-biocube-bydgoszcz-2024")
	if item.URL != "New URL" {
		t.Error("Expected change in the URL, got ", item.URL)
	}

	item = n.Get("cube4fun-biocube-bydgoszcz-2024")
	if !item.HasPassed {
		t.Error("Expected hasPassed = true, got ", item.HasPassed)
	}
}
