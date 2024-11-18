package database

import (
	"database/sql"
	"github.com/hrabit64/springnote-breezenote/config"
	"github.com/hrabit64/springnote-breezenote/pkg/utils"
	_ "github.com/mattn/go-sqlite3"
	"path"
)

// GetConnect DB 커넥션을 가져옵니다.
func GetConnect() (*sql.DB, error) {

	conn, err := sql.Open("sqlite3", path.Join(utils.GetRootPath(), config.RootConfig.DBConnURL))

	if err != nil {
		return nil, err
	}

	return conn, err
}
