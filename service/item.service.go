package service

import (
	"database/sql"
	"github.com/hrabit64/springnote-breezenote/database"
)

type ItemService interface {
	ReadItemById(id string) (*database.Item, error)
	CreateItem(item *database.Item) error
}

type itemService struct {
}

func NewItemService() ItemService {
	return &itemService{}
}

// ReadItemById 아이템을 아이디로 조회합니다.
func (i *itemService) ReadItemById(id string) (*database.Item, error) {
	conn, err := database.GetConnect()
	if err != nil {
		return nil, err
	}

	defer conn.Close()

	item, err := getItemById(conn, id)
	if err != nil {
		return nil, err
	}

	return item, nil
}

// CreateItem 아이템을 생성합니다.
func (i *itemService) CreateItem(item *database.Item) error {
	conn, err := database.GetConnect()
	if err != nil {
		return err
	}

	defer conn.Close()

	err = createItem(conn, item)
	if err != nil {
		return err
	}

	return nil
}

func getItemById(conn *sql.DB, id string) (*database.Item, error) {
	query := `SELECT ITEM_PK, ORIGIN_NAME, CREATE_AT FROM ITEM WHERE ITEM_PK = ?`

	stmt, err := conn.Prepare(query)
	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	var item database.Item
	err = stmt.QueryRow(id).Scan(&item.Id, &item.OriginName, &item.CreateAt)
	if err != nil {
		return nil, err
	}

	return &item, nil
}

func createItem(conn *sql.DB, item *database.Item) error {
	query := `INSERT INTO ITEM (ITEM_PK, ORIGIN_NAME, CREATE_AT) VALUES (?, ?, ?)`
	stmt, err := conn.Prepare(query)
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(item.Id, item.OriginName, item.CreateAt)
	if err != nil {
		return err
	}

	return nil

}
