package main

import (
	"database/sql"
)

func ReadRepos(db sqlQueryable, orderby string, limit, offset int) ([]*Repo, error) {
	rows, err := db.Query(qRepos, "created DESC", limit, offset)
	if err != nil {
		return nil, err
	}

	return unmarshalRepos(rows, limit)
}

func unmarshalRepos(rows *sql.Rows, limit int) ([]*Repo, error) {
	defer rows.Close()
	repos := make([]*Repo, limit)
	i := 0
	for rows.Next() {
		r := &Repo{}
		if err := r.UnmarshalSQL(rows); err != nil {
			return nil, err
		}
		repos[i] = r
		i++
	}

	repos = repos[:i]
	return repos, nil
}
