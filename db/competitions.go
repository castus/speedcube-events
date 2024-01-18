package db

type Competitions []Competition

func (c Competitions) FindByID(ID string) *Competition {
	for _, item := range c {
		if item.Id == ID {
			return &item
		}
	}

	return nil
}
