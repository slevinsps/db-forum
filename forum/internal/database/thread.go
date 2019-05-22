package database

import (
	"database/sql"
	"db_forum/internal/models"
	"db_forum/internal/utils"
	"strconv"
)

func (db DataBase) isThreadUnique(thread models.Thread) (threadRes models.Thread, checkUnique bool, err error) {
	var tx *sql.Tx
	tx, err = db.Db.Begin()
	defer tx.Rollback()
	checkUnique = false

	sqlStatement := "SELECT t.author, t.created, t.forum, t.id, t.message, t.title, t.slug, t.votes " +
		"FROM Thread as t where  t.slug = $1"
	row := tx.QueryRow(sqlStatement, thread.Slug)
	err = row.Scan(&threadRes.Author, &threadRes.Created, &threadRes.Forum, &threadRes.Id, &threadRes.Message, &threadRes.Title, &threadRes.Slug, &threadRes.Votes)

	if err != nil {
		if err == sql.ErrNoRows {
			checkUnique = true
			err = nil
		}
		utils.PrintDebug("database/isThreadUnique Query error")

		return
	}

	err = tx.Commit()
	if err != nil {
		return
	}

	return
}

// CreateThread
func (db *DataBase) CreateThread(thread models.Thread) (threadQuery models.Thread, checkUnique bool, err error) {
	var tx *sql.Tx
	tx, err = db.Db.Begin()
	defer tx.Rollback()

	checkUnique = true
	if thread.Slug != "" {
		if threadQuery, checkUnique, err = db.isThreadUnique(thread); err != nil {
			utils.PrintDebug("database/CreateThread - fail uniqie")
			return
		}
	}

	if !checkUnique {
		utils.PrintDebug("CreateThread ", thread)
		return
	}

	if thread.Created == "" {
		sqlInsert := `
		INSERT INTO Thread(author, forum, message, title, slug) VALUES
			($1, $2, $3, $4, $5) RETURNING author, forum, id, message, title, slug;
			`
		row := tx.QueryRow(sqlInsert, thread.Author, thread.Forum, thread.Message, thread.Title, thread.Slug)
		err = row.Scan(&threadQuery.Author, &threadQuery.Forum, &threadQuery.Id, &threadQuery.Message, &threadQuery.Title, &threadQuery.Slug)
	} else {
		sqlInsert := `
		INSERT INTO Thread(author, created, forum, message, title, slug) VALUES
			($1, $2, $3, $4, $5, $6) RETURNING author, created, forum, id, message, title, slug;
			`
		row := tx.QueryRow(sqlInsert, thread.Author, thread.Created, thread.Forum, thread.Message, thread.Title, thread.Slug)
		err = row.Scan(&threadQuery.Author, &threadQuery.Created, &threadQuery.Forum, &threadQuery.Id, &threadQuery.Message, &threadQuery.Title, &threadQuery.Slug)
	}
	if err != nil {
		return
	}
	err = tx.Commit()
	if err != nil {
		return
	}
	utils.PrintDebug("database/CreateThread +")

	return
}

// GetThreadsByForum
func (db *DataBase) GetThreadsByForum(title string, limitStr string, sinceStr string, descStr string) (threads []models.Thread, err error) {
	var tx *sql.Tx
	tx, err = db.Db.Begin()
	defer tx.Rollback()

	sqlQuery :=
		"SELECT t.author, t.created, t.forum, t.id, t.message, t.slug, t.title, t.votes " +
			"FROM Thread as t where t.forum = $1 "
	if sinceStr != "" {
		if descStr == "true" {
			sqlQuery += "and t.created <= " + "'" + sinceStr + "'"
		} else {
			sqlQuery += "and t.created >= " + "'" + sinceStr + "'"
		}
	}
	sqlQuery += " order by t.created "
	if descStr == "true" {
		sqlQuery += "desc "
	}
	if limitStr != "" {
		sqlQuery += "limit " + limitStr + ";"
	}

	rows, erro := tx.Query(sqlQuery, title)
	if erro != nil {
		err = erro
		utils.PrintDebug("database/GetThreadsByForum Query error")
		return
	}

	defer rows.Close()

	for rows.Next() {
		thread := models.Thread{}
		if err = rows.Scan(&thread.Author, &thread.Created, &thread.Forum, &thread.Id, &thread.Message, &thread.Slug, &thread.Title, &thread.Votes); err != nil {
			utils.PrintDebug("database/GetThreadsByForum wrong row catched")
			break
		}
		threads = append(threads, thread)
	}

	err = tx.Commit()
	if err != nil {
		return
	}

	utils.PrintDebug("database/GetThreadsByForum +")

	return
}

//var maxElapsed time.Duration = -1

// GetThreadsById
func (db *DataBase) GetThreadById(slugOrId string) (thread models.Thread, checkFindThread bool, err error) {

	checkFindThread = true

	slugOrIDInt, errAtoi := strconv.Atoi(slugOrId)
	//start := time.Now()
	if errAtoi != nil {

		sqlQuery :=
			"SELECT t.author, t.created, t.forum, t.id, t.message, t.slug, t.title, t.votes " +
				"FROM Thread as t where t.slug = $1;"
		row := db.Db.QueryRow(sqlQuery, slugOrId)
		err = row.Scan(&thread.Author, &thread.Created, &thread.Forum, &thread.Id, &thread.Message, &thread.Slug, &thread.Title, &thread.Votes)
	} else {
		sqlQuery :=
			"SELECT t.author, t.created, t.forum, t.id, t.message, t.slug, t.title, t.votes " +
				"FROM Thread as t where t.id = $1;"
		row := db.Db.QueryRow(sqlQuery, slugOrIDInt)
		err = row.Scan(&thread.Author, &thread.Created, &thread.Forum, &thread.Id, &thread.Message, &thread.Slug, &thread.Title, &thread.Votes)
	}
	/*end := time.Now()
	elapsed := end.Sub(start)
	if maxElapsed < elapsed {
		maxElapsed = elapsed
		fmt.Println("strconv.Atoi(slugOrId) duration (", slugOrId, ") ", elapsed.String())
	}*/

	if err != nil {
		if err == sql.ErrNoRows {
			checkFindThread = false
			err = nil
		}
		utils.PrintDebug("database/GetThreadsById Scan error")
		return
	}

	utils.PrintDebug("database/GetThreadsById +")

	return
}

// UpdateThread
func (db *DataBase) UpdateThread(threadNew models.Thread, threadOld models.Thread) (threadRes models.Thread, err error) {

	var tx *sql.Tx
	tx, err = db.Db.Begin()
	defer tx.Rollback()

	//utils.PrintDebug(user)
	sqlInsert := `
		UPDATE Thread SET message = $1, title = $2 where id = $3 RETURNING author, created, forum, id, message, title, slug, votes;
		`
	if threadNew.Message == "" {
		threadNew.Message = threadOld.Message
	}
	if threadNew.Title == "" {
		threadNew.Title = threadOld.Title
	}
	row := tx.QueryRow(sqlInsert, threadNew.Message, threadNew.Title, threadOld.Id)
	err = row.Scan(&threadRes.Author, &threadRes.Created, &threadRes.Forum, &threadRes.Id, &threadRes.Message, &threadRes.Title, &threadRes.Slug, &threadRes.Votes)
	if err != nil {
		return
	}
	err = tx.Commit()
	if err != nil {
		return
	}
	utils.PrintDebug("database/UpdateThread +")

	return
}

// CountThread
func (db *DataBase) CountThread() (count int, err error) {
	var tx *sql.Tx
	tx, err = db.Db.Begin()
	defer tx.Rollback()

	//utils.PrintDebug(user)
	sqlInsert := `
		SELECT COUNT(*) FROM Thread
		`
	row := tx.QueryRow(sqlInsert)
	err = row.Scan(&count)
	if err != nil {
		return
	}
	err = tx.Commit()
	if err != nil {
		return
	}
	utils.PrintDebug("database/CountThread +")

	return
}
