package database

import (
	"database/sql"
)

// InitSchema 스키마를 초기화합니다.
func InitSchema(conn *sql.DB) error {

	_, err := conn.Exec(`
		CREATE TABLE IF NOT EXISTS ITEM (
			ITEM_PK TEXT PRIMARY KEY,
			ORIGIN_NAME TEXT NOT NULL,
			CREATE_AT TIMESTAMP NOT NULL
		);`)

	if err != nil {
		return err
	}

	return nil

}

func RunSetup() error {

	conn, err := GetConnect()
	if err != nil {
		return err
	}
	defer conn.Close()

	err = InitSchema(conn)
	if err != nil {
		return err
	}

	return nil

}
