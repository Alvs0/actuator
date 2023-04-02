package impl

import (
	"actuator/engine"
)

type AccountQuery interface {
	GetUser(username string) ([]UserDb, error)
}

type accountQuery struct {
	sqlAdapter engine.SqlAdapter
}

func NewAccountQuery(sqlAdapter engine.SqlAdapter) AccountQuery {
	return &accountQuery{
		sqlAdapter: sqlAdapter,
	}
}

type UserDb struct {
	ID                string `json:"id" db:"id"`
	Username          string `json:"username" db:"username"`
	EncryptedPassword string `json:"password" db:"password"`
}

func (a *accountQuery) GetUser(username string) ([]UserDb, error) {
	getQuery := `SELECT id, username, password FROM account WHERE username = ?`

	queryInput := []interface{}{
		username,
	}

	var userDbs []UserDb
	if err := a.sqlAdapter.Read(getQuery, queryInput, &userDbs); err != nil {
		return nil, err
	}

	return userDbs, nil
}
