package storage

import (
	"database/sql"
	"errors"
	_ "github.com/mattn/go-sqlite3"
	"time"
)

type Storage interface {
	GetGithubAuth(user string) (string, error)
	SaveGithubAuth(user, auth string) error
	GetRepos(needle string) ([]*Repo, error)
	GetRepo(needle string) (*Repo, error)
	AddRepo(myType string, origin string, last time.Time, ready bool) error
}

type Repo struct {
	Type   *string
	Origin *string
	Ready  *bool
	Last   *time.Time
}

type repoImpl struct {
	connection *sql.DB
	noRepos    []*Repo
}

func (r repoImpl) SaveGithubAuth(user, auth string) error {
	stmt, err := r.connection.Prepare("insert or replace into github_auth (login, auth) values(?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Query(user, auth)
	return err
}

func (r repoImpl) GetGithubAuth(user string) (string, error) {
	stmt, err := r.connection.Prepare("select auth from github_auth where login = ?")
	if err != nil {
		return "", err
	}
	defer stmt.Close()
	rows, err := stmt.Query(user)
	if err != nil {
		return "", err
	}
	defer rows.Close()
	for rows.Next() {
		var auth string
		rows.Scan(&auth)
		return auth, nil
	}
	return "", errors.New("No auth found for " + user)
}

func (r repoImpl) AddRepo(myType string, origin string, last time.Time, ready bool) error {
	stmt, err := r.connection.Prepare("delete from repositories where type = ? and origin = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(myType, origin)
	if err != nil {
		return err
	}
	stmt, err = r.connection.Prepare("insert into repositories (type, origin, last, ready) values(?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(myType, origin, last, ready)
	return err
}

func (r repoImpl) GetRepo(needle string) (*Repo, error) {
	stmt, err := r.connection.Prepare("select origin, last, ready from repositories where origin = ? and type = 'github'")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	rows, err := stmt.Query(needle)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		myType := "github"
		var last time.Time
		var origin string
		var ready bool
		rows.Scan(&origin, &last, &ready)
		return &Repo{Last: &last, Origin: &origin, Ready: &ready, Type: &myType}, nil
	}
	return nil, nil
}

func (r repoImpl) GetRepos(needle string) ([]*Repo, error) {
	stmt, err := r.connection.Prepare("select type, origin, last, ready from repositories where id like ?")
	if err != nil {
		return r.noRepos, err
	}
	defer stmt.Close()
	rows, err := stmt.Query("%" + needle + "%")
	if err != nil {
		return r.noRepos, err
	}
	defer rows.Close()
	result := make([]*Repo, 0, 20)
	for rows.Next() {
		var last time.Time
		var origin string
		var myType string
		var ready bool
		rows.Scan(&myType, &origin, &last, &ready)
		result = append(result, &Repo{Origin: &origin, Last: &last, Type: &myType, Ready: &ready})
	}
	return result, nil
}

func Create(filename string) (*Storage, error) {
	db, err := sql.Open("sqlite3", filename)
	if err != nil {
		return nil, err
	}
	sqlStmt := `
		  create table if not exists repositories 
			(
				id integer not null primary key autoincrement, 
				type text not null, 
				origin text not null,
				ready integer not null,
				last timestamp not null
			) ;
			create table if not exists github_auth
			(
				login text primary key,
				auth text
			);
			`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		return nil, err
	}
	var result Storage = Storage(repoImpl{connection: db, noRepos: []*Repo{}})
	return &result, nil
}
