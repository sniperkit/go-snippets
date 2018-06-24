package tasks

import (
	"database/sql"
	"fmt"
	// "github.com/datatogether/sql_datastore"
	"github.com/ipfs/go-datastore"
	"github.com/ipfs/go-datastore/query"
)

// ReadTasks reads a list of tasks from store
func ReadTasks(store datastore.Datastore, orderby string, limit, offset int) ([]*Task, error) {
	q := query.Query{
		Prefix: fmt.Sprintf("/%s", Task{}.DatastoreType()),
		Limit:  limit,
		Offset: offset,
		// TODO - add native ordering support
		// Orders: []query.Order{}
	}

	res, err := store.Query(q)
	if err != nil {
		return nil, err
	}

	tasks := make([]*Task, limit)
	i := 0
	for r := range res.Next() {
		if r.Error != nil {
			return nil, err
		}

		c, ok := r.Value.(*Task)
		if !ok {
			return nil, fmt.Errorf("Invalid Response")
		}

		tasks[i] = c
		i++
	}

	return tasks[:i], nil
}

// TODO - transfer to kiwix taskdef
// func GenerateAvailableTasks(db *sql.DB) ([]*Task, error) {
// 	row, err := db.Query(qAvailableTasks)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// TODO - make not bad
// 	store := sql_datastore.NewDatastore(db)
// 	store.Register(&Task{})

// 	tasks := []*Task{}
// 	for row.Next() {
// 		var (
// 			repoUrl, repoCommit, sourceTitle, sourceUrl, sourceChecksum string
// 		)
// 		if err := row.Scan(&repoUrl, &repoCommit, &sourceTitle, &sourceUrl, &sourceChecksum); err != nil {
// 			return nil, err
// 		}

// 		t := &Task{
// 			Title:          fmt.Sprintf("injest %s to ipfs", sourceTitle),
// 			RepoUrl:        repoUrl,
// 			RepoCommit:     repoCommit,
// 			SourceUrl:      sourceUrl,
// 			SourceChecksum: sourceChecksum,
// 		}

// 		if err := t.Save(store); err != nil {
// 			return nil, err
// 		}

// 		tasks = append(tasks, t)
// 	}

// 	return tasks, nil
// }

func unmarshalTasks(rows *sql.Rows, limit int) ([]*Task, error) {
	defer rows.Close()
	tasks := make([]*Task, limit)
	i := 0
	for rows.Next() {
		t := &Task{}
		if err := t.UnmarshalSQL(rows); err != nil {
			return nil, err
		}
		tasks[i] = t
		i++
	}

	tasks = tasks[:i]
	return tasks, nil
}
