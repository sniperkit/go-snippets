package main

import "time"

type Calendar struct {
	Members      []string  `json:"members" bson:"members"`
	Moments      []Moment  `json:"moments" bson:"moments"`
	Created      time.Time `json:"created" bson:"created"`
	LastModified time.Time `json:"lastModified" bson:"lastModified"`
}

func (c *Calendar) Within(start, end string) ([]Moment, error) {
	s, err := time.Parse("2006-01-02", start)
	if err != nil {
		return nil, err
	}

	e, err := time.Parse("2006-01-02", end)
	if err != nil {
		return nil, err
	}

	moments := []Moment{}
	for _, m := range c.Moments {
		if (s.Before(m.StartDate) || s.Equal(m.StartDate)) || (e.After(m.EndDate) || e.Equal(m.EndDate)) {
			moments = append(moments, m)
		} else if m.Repeating {
			// TODO: Work in repeating events
		}
	}

	return moments, nil
}
