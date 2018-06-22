package repo

import (
	r "github.com/dancannon/gorethink"

	"github.com/gophergala/Gevvent/services/gevvent-event-service/model"
	"github.com/gophergala/Gevvent/services/gevvent-lib/rethinkdb"
)

func CreateEvent(event *model.Event) (*model.Event, error) {
	s, err := rethinkdb.Session()
	if err != nil {
		return nil, err
	}
	response, err := r.Table("events").Insert(event).RunWrite(s)
	if err != nil {
		return nil, err
	}

	// Find new ID of product if needed
	if event.ID == "" && len(response.GeneratedKeys) == 1 {
		event.ID = response.GeneratedKeys[0]
	}

	return event, nil
}

func DeleteEvent(id string) error {
	s, err := rethinkdb.Session()
	if err != nil {
		return err
	}

	return r.Table("events").Get(id).Delete().Exec(s)
}

func Get(id string) (*model.Event, error) {
	var event *model.Event

	s, err := rethinkdb.Session()
	if err != nil {
		return nil, err
	}

	rows, err := r.Table("events").Get(id).Run(s)
	if err != nil {
		return nil, err
	}

	err = rows.One(&event)
	if err != nil {
		if err == r.ErrEmptyResult {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return event, err
}

type SearchResult struct {
	*model.Event

	Distance float64 `gorethink:"distance", json:"distance"`
}

func Search(userID, query string, page, count int64) ([]*SearchResult, error) {
	var events []*SearchResult

	s, err := rethinkdb.Session()
	if err != nil {
		return events, err
	}

	rows, err := r.Table("events").Filter(func(row r.Term) r.Term {
		return row.Field("name").Match(query).And(
			r.Not(row.Field("private")).Or(row.Field("user_id").Eq(userID)),
		)
	}).Skip((page - 1) * count).Limit(count).Run(s)
	if err != nil {
		return events, err
	}

	err = rows.All(&events)
	if err != nil {
		return events, err
	}

	return events, err
}

func SearchNearest(userID, query string, lat, lng float64, page, count int64) ([]*SearchResult, error) {
	var events []*SearchResult

	s, err := rethinkdb.Session()
	if err != nil {
		return events, err
	}

	rows, err := r.Table("events").GetNearest(
		r.Point(lng, lat),
		r.GetNearestOpts{
			Index: "whereLocation",
			Unit:  "km",
		},
	).Map(func(row r.Term) interface{} {
		return row.Field("doc").Merge(map[string]interface{}{
			"distance": row.Field("dist"),
		})
	}).Filter(func(row r.Term) r.Term {
		return row.Field("name").Match(query).And(
			r.Not(row.Field("private")).Or(row.Field("user_id").Eq(userID)),
		)
	}).Skip((page - 1) * count).Limit(count).Run(s)
	if err != nil {
		return events, err
	}

	err = rows.All(&events)
	if err != nil {
		return events, err
	}

	return events, err
}

func GetNewest(page, count int64) ([]*model.Event, error) {
	var events []*model.Event

	s, err := rethinkdb.Session()
	if err != nil {
		return events, err
	}

	rows, err := r.Table("events").OrderBy(r.Row.Field("when").Field("from")).
		Filter(r.Not(r.Row.Field("private"))).
		Filter(r.Row.Field("when").Field("from").Ge(r.Now())).
		Skip((page - 1) * count).Limit(count).Run(s)
	if err != nil {
		return events, err
	}

	err = rows.All(&events)
	if err != nil {
		return events, err
	}

	return events, err
}

func GetHosting(userID string, page, count int64) ([]*model.Event, error) {
	var events []*model.Event

	s, err := rethinkdb.Session()
	if err != nil {
		return events, err
	}

	rows, err := r.Table("events").Filter(func(row r.Term) r.Term {
		return row.Field("user_id").Match(userID)
	}).Skip((page - 1) * count).Limit(count).Run(s)
	if err != nil {
		return events, err
	}

	err = rows.All(&events)
	if err != nil {
		return events, err
	}

	return events, err
}
