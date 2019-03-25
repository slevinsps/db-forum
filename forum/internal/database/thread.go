package database

import (
	"database/sql"
	"db_forum/internal/models"
	"fmt"
	"strconv"
)

func (db DataBase) isThreadUnique(thread models.Thread, checkUnique *bool) (threadRes models.Thread, err error) {
	var tx *sql.Tx
	tx, err = db.Db.Begin()
	defer tx.Rollback()
	*checkUnique = false

	if thread.Created != "" {
		if thread.Slug != "" {
			sqlStatement := "SELECT t.author, t.created, t.forum, t.id, t.message, t.title, t.slug " +
				"FROM Thread as t where  lower(t.slug) like lower($1)"
			row := tx.QueryRow(sqlStatement, thread.Slug)
			err = row.Scan(&threadRes.Author, &threadRes.Created, &threadRes.Forum, &threadRes.Id, &threadRes.Message, &threadRes.Title, &threadRes.Slug)
		} else {
			sqlStatement := "SELECT t.author, t.created, t.forum, t.id, t.message, t.title " +
				"FROM Thread as t where  lower(t.title) like lower($1)"
			row := tx.QueryRow(sqlStatement, thread.Title)
			err = row.Scan(&threadRes.Author, &threadRes.Created, &threadRes.Forum, &threadRes.Id, &threadRes.Message, &threadRes.Title)
		}
	} else {
		if thread.Slug != "" {
			sqlStatement := "SELECT t.author, t.forum, t.id, t.message, t.title, t.slug " +
				"FROM Thread as t where  lower(t.slug) like lower($1)"
			row := tx.QueryRow(sqlStatement, thread.Slug)
			err = row.Scan(&threadRes.Author, &threadRes.Forum, &threadRes.Id, &threadRes.Message, &threadRes.Title, &threadRes.Slug)
		} else {
			sqlStatement := "SELECT t.author,  t.forum, t.id, t.message, t.title " +
				"FROM Thread as t where  lower(t.title) like lower($1)"
			row := tx.QueryRow(sqlStatement, thread.Title)
			err = row.Scan(&threadRes.Author, &threadRes.Forum, &threadRes.Id, &threadRes.Message, &threadRes.Title)
		}
	}

	if err != nil {
		if err == sql.ErrNoRows {
			*checkUnique = true
			err = nil
		}
		fmt.Println("database/isThreadUnique Query error")
		return
	}

	err = tx.Commit()
	if err != nil {
		return
	}

	return
}

// CreateThread
func (db *DataBase) CreateThread(thread models.Thread, slug string) (threadQuery models.Thread, checkUnique bool, err error) {
	var tx *sql.Tx
	tx, err = db.Db.Begin()
	defer tx.Rollback()

	checkUnique = false
	if threadQuery, err = db.isThreadUnique(thread, &checkUnique); err != nil {
		fmt.Println("database/CreateUser - fail uniqie")
		return
	}

	if !checkUnique {
		fmt.Println("CreateThread ", thread)
		return
	}

	//fmt.Println(user)
	if thread.Created != "" {
		if thread.Slug != "" {
			sqlInsert := `
			INSERT INTO Thread(author, created, forum, message, title, slug) VALUES
				($1, $2, $3, $4, $5, $6) RETURNING author, created, forum, id, message, title, slug;
				`
			row := tx.QueryRow(sqlInsert, thread.Author, thread.Created, thread.Forum, thread.Message, thread.Title, thread.Slug)
			err = row.Scan(&threadQuery.Author, &threadQuery.Created, &threadQuery.Forum, &threadQuery.Id, &threadQuery.Message, &threadQuery.Title, &threadQuery.Slug)
			if err != nil {
				return
			}
		} else {
			sqlInsert := `
			INSERT INTO Thread(author, created, forum, message, title) VALUES
				($1, $2, $3, $4, $5) RETURNING author, created, forum, id, message, title;
				`
			row := tx.QueryRow(sqlInsert, thread.Author, thread.Created, thread.Forum, thread.Message, thread.Title)
			err = row.Scan(&threadQuery.Author, &threadQuery.Created, &threadQuery.Forum, &threadQuery.Id, &threadQuery.Message, &threadQuery.Title)
			if err != nil {
				return
			}
		}
	} else {
		if thread.Slug != "" {
			sqlInsert := `
			INSERT INTO Thread(author, forum, message, title, slug) VALUES
				($1, $2, $3, $4, $5) RETURNING author, forum, id, message, title, slug;
				`
			row := tx.QueryRow(sqlInsert, thread.Author, thread.Forum, thread.Message, thread.Title, thread.Slug)
			err = row.Scan(&threadQuery.Author, &threadQuery.Forum, &threadQuery.Id, &threadQuery.Message, &threadQuery.Title, &threadQuery.Slug)
			if err != nil {
				return
			}
		} else {
			sqlInsert := `
			INSERT INTO Thread(author, forum, message, title) VALUES
				($1, $2, $3, $4) RETURNING author, forum, id, message, title;
				`
			row := tx.QueryRow(sqlInsert, thread.Author, thread.Forum, thread.Message, thread.Title)
			err = row.Scan(&threadQuery.Author, &threadQuery.Forum, &threadQuery.Id, &threadQuery.Message, &threadQuery.Title)
			if err != nil {
				return
			}
		}
	}

	_ = db.UpdateFieldsForum(thread.Forum, 1, "threads")
	if err != nil {
		return
	}
	err = tx.Commit()
	if err != nil {
		return
	}
	fmt.Println("database/CreateThread +")

	return
}

// GetThreadsByForum
func (db *DataBase) GetThreadsByForum(title string, limitStr string, sinceStr string, descStr string) (threads []models.Thread, err error) {
	var tx *sql.Tx
	tx, err = db.Db.Begin()
	defer tx.Rollback()

	sqlQuery :=
		"SELECT t.author, t.created, t.forum, t.id, t.message, t.slug, t.title " +
			"FROM Thread as t where lower(t.forum) like lower($1) "
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
		fmt.Println("database/GetThreadsByForum Query error")
		return
	}

	defer rows.Close()

	for rows.Next() {
		thread := models.Thread{}
		if err = rows.Scan(&thread.Author, &thread.Created, &thread.Forum, &thread.Id, &thread.Message, &thread.Slug, &thread.Title); err != nil {
			fmt.Println("database/GetThreadsByForum wrong row catched")
			break
		}
		threads = append(threads, thread)
	}

	err = tx.Commit()
	if err != nil {
		return
	}

	fmt.Println("database/GetThreadsByForum +")

	return
}

// GetThreadsById
func (db *DataBase) GetThreadById(slugOrId string) (thread models.Thread, checkFindThread bool, err error) {
	var tx *sql.Tx
	tx, err = db.Db.Begin()
	defer tx.Rollback()
	checkFindThread = true
	slugOrIDInt, errAtoi := strconv.Atoi(slugOrId)
	if errAtoi != nil {

		sqlQuery :=
			"SELECT t.author, t.created, t.forum, t.id, t.message, t.slug, t.title, t.votes " +
				"FROM Thread as t where lower(t.slug) like lower($1);"
		row := tx.QueryRow(sqlQuery, slugOrId)
		err = row.Scan(&thread.Author, &thread.Created, &thread.Forum, &thread.Id, &thread.Message, &thread.Slug, &thread.Title, &thread.Votes)
	} else {
		sqlQuery :=
			"SELECT t.author, t.created, t.forum, t.id, t.message, t.slug, t.title, t.votes " +
				"FROM Thread as t where t.id = $1;"
		row := tx.QueryRow(sqlQuery, slugOrIDInt)
		err = row.Scan(&thread.Author, &thread.Created, &thread.Forum, &thread.Id, &thread.Message, &thread.Slug, &thread.Title, &thread.Votes)
	}

	if err != nil {
		if err == sql.ErrNoRows {
			checkFindThread = false
			err = nil
		}
		fmt.Println("database/GetThreadsById Scan error")
		return
	}

	err = tx.Commit()
	if err != nil {
		fmt.Println("database/GetThreadsById Commit error")
		return
	}

	fmt.Println("database/GetThreadsById +")

	return
}

// UpdateThread
func (db *DataBase) UpdateThread(threadNew models.Thread, threadOld models.Thread) (threadRes models.Thread, err error) {

	var tx *sql.Tx
	tx, err = db.Db.Begin()
	defer tx.Rollback()

	//fmt.Println(user)
	sqlInsert := `
		UPDATE Thread SET message = $1, title = $2 where id = $3 RETURNING author, created, forum, id, message, title, slug;
		`
	if threadNew.Message == "" {
		threadNew.Message = threadOld.Message
	}
	if threadNew.Title == "" {
		threadNew.Title = threadOld.Title
	}
	row := tx.QueryRow(sqlInsert, threadNew.Message, threadNew.Title, threadOld.Id)
	err = row.Scan(&threadRes.Author, &threadRes.Created, &threadRes.Forum, &threadRes.Id, &threadRes.Message, &threadRes.Title, &threadRes.Slug)
	if err != nil {
		return
	}
	err = tx.Commit()
	if err != nil {
		return
	}
	fmt.Println("database/UpdateThread +")

	return
}

// CountThread
func (db *DataBase) CountThread() (count int, err error) {
	var tx *sql.Tx
	tx, err = db.Db.Begin()
	defer tx.Rollback()

	//fmt.Println(user)
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
	fmt.Println("database/CountThread +")

	return
}
