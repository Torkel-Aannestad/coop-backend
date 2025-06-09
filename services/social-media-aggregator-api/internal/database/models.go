package database

import "database/sql"

type Models struct {
	Messages *MessageModel
}

func NewModels(db *sql.DB) *Models {
	return &Models{
		Messages: &MessageModel{DB: db},
	}
}
