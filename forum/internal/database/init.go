package database

import (
	"database/sql"
	"db_forum/internal/config"
	"db_forum/internal/utils"
	"io/ioutil"
	"os"

	//
	_ "github.com/lib/pq"
)

// Init try to connect to DataBase.
// If success - return instance of DataBase
// if failed - return error
func Init(CDB config.DatabaseConfig) (db *DataBase, err error) {

	// for local launch
	if os.Getenv(CDB.URL) == "" {
		os.Setenv(CDB.URL, "user=db_forum_user password=db_forum_password dbname=db_forum sslmode=disable")
	}

	var database *sql.DB
	if database, err = sql.Open(CDB.DriverName, os.Getenv(CDB.URL)); err != nil {
		utils.PrintDebug("database/Init cant open:" + err.Error())
		return
	}

	db = &DataBase{
		Db: database,
	}
	db.Db.SetMaxOpenConns(CDB.MaxOpenConns)

	if err = db.Db.Ping(); err != nil {
		utils.PrintDebug("database/Init cant access:" + err.Error())
		return
	}
	utils.PrintDebug("database/Init open")
	if err = db.CreateTables(); err != nil {
		return
	}
	return
}

// CountForum
func (db *DataBase) ServiceClear() (err error) {
	var tx *sql.Tx
	tx, err = db.Db.Begin()
	defer tx.Rollback()

	//utils.PrintDebug(user)
	sqlRow := `
	  TRUNCATE Post, Users, Thread, Forum, Vote;
		`
	_, err = tx.Exec(sqlRow)
	if err != nil {
		return
	}
	err = tx.Commit()
	if err != nil {
		return
	}
	utils.PrintDebug("database/ServiceClear +")
	return
}

func (db *DataBase) CreateTables() error {
	query, err := ioutil.ReadFile("init.pgsql")
	if err != nil {
		panic(err)
	}
	_, err = db.Db.Exec(string(query))
	if err != nil {
		utils.PrintDebug("database/init - fail:" + err.Error())
	}
	return err
}
