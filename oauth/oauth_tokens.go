package oauth

import (
	"database/sql"
)

func UnmarshalTokens(rows *sql.Rows) ([]*UserOauthToken, error) {
	tokens := []*UserOauthToken{}
	defer rows.Close()
	for rows.Next() {
		t := &UserOauthToken{}
		if err := t.UnmarshalSQL(rows); err != nil {
			return nil, err
		}
		tokens = append(tokens, t)
	}
	return tokens, nil
}
