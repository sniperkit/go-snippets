package dshash

type PersonsService struct {
	PersonRepository
}

func (l *PersonsService) GetAll(person *Person) ([]string, error) {
	foundPerson, e := l.PersonRepository.Find(person)

	if e != nil {
		return nil, e
	}

	if foundPerson == nil {
		return []string{}, nil
	} else {
		return foundPerson.Locations, nil
	}
}

func (l *PersonsService) Save(person *Person) error {
	return l.PersonRepository.Save(person)
}
