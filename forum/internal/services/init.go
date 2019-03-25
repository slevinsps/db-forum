package api

import (
	"db_forum/internal/config"
	"db_forum/internal/database"
)

// creates Handler
func Init(DB *database.DataBase) (handler *Handler) {
	handler = &Handler{
		DB: *DB,
	}
	return
}

func GetHandler(confPath string) (handler *Handler, conf *config.Configuration, err error) {

	var (
		db *database.DataBase
	)

	if conf, err = config.Init(confPath); err != nil {
		return
	}

	if db, err = database.Init(conf.DataBase); err != nil {
		return
	}

	handler = Init(db)
	return
}
