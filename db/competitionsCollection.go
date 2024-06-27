package db

type CompetitionsCollection []*Competition

func (c *CompetitionsCollection) FilterNotOnline() CompetitionsCollection {
	var items = CompetitionsCollection{}
	for _, competition := range *c {
		if competition.Place != "zawody online" {
			items = append(items, competition)
		}
	}

	return items
}

func (c *CompetitionsCollection) FilterNotPassed() CompetitionsCollection {
	var items = CompetitionsCollection{}
	for _, competition := range *c {
		if !competition.HasPassed {
			items = append(items, competition)
		}
	}

	return items
}

func (c *CompetitionsCollection) FilterEmptyEvents() CompetitionsCollection {
	var items = CompetitionsCollection{}
	for _, competition := range *c {
		if len(competition.Events) == 0 {
			items = append(items, competition)
		}
	}

	return items
}

func (c *CompetitionsCollection) FilterEmptyMainEvent() CompetitionsCollection {
	var items = CompetitionsCollection{}
	for _, competition := range *c {
		if competition.MainEvent == "" {
			items = append(items, competition)
		}
	}

	return items
}

func (c *CompetitionsCollection) FilterEmptyCompetitorLimit() CompetitionsCollection {
	var items = CompetitionsCollection{}
	for _, competition := range *c {
		if competition.CompetitorLimit == 0 {
			items = append(items, competition)
		}
	}

	return items
}

func (c *CompetitionsCollection) FilterEmptyRegistered() CompetitionsCollection {
	var items = CompetitionsCollection{}
	for _, competition := range *c {
		if competition.Registered == 0 {
			items = append(items, competition)
		}
	}

	return items
}

func (c *CompetitionsCollection) FilterWCAEvents() CompetitionsCollection {
	var items = CompetitionsCollection{}
	for _, competition := range *c {
		if competition.Type == CompetitionType.WCA {
			items = append(items, competition)
		}
	}

	return items
}

func (c *CompetitionsCollection) FilterEmptyDistanceOrDuration() CompetitionsCollection {
	var items = CompetitionsCollection{}
	for _, competition := range *c {
		if competition.Distance == "" && competition.Duration == "" {
			items = append(items, competition)
		}
	}

	return items
}
