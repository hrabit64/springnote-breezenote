package database

import "time"

type Item struct {
	Id         string    // ITEM_PK
	OriginName string    // ORIGIN_NAME
	CreateAt   time.Time // CREATE_AT
}
