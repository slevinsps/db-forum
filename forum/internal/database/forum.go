package database

import (
	"database/sql"
	"db_forum/internal/models"
	"fmt"
)

// GetForumBySlug
func (db *DataBase) GetForumBySlug(slug string) (forum models.Forum, checkFindForum bool, err error) {
	checkFindForum = true
	var tx *sql.Tx
	tx, err = db.Db.Begin()
	defer tx.Rollback()

	sqlQuery :=
		"SELECT f.posts, f.slug, f.threads, f.title, f.user FROM Forum as f " +
			"where lower(f.slug) like lower($1);"

	row := tx.QueryRow(sqlQuery, slug)
	err = row.Scan(&forum.Posts, &forum.Slug, &forum.Threads, &forum.Title, &forum.User)
	if err != nil {
		if err == sql.ErrNoRows {
			checkFindForum = false
			err = nil
		}
		return
	}
	err = tx.Commit()
	if err != nil {
		return
	}

	fmt.Println("database/GetForumBySlug +")

	return
}

// CreateThread
func (db *DataBase) UpdateFieldsForum(slug string, number int, field string) (err error) {
	var tx *sql.Tx
	tx, err = db.Db.Begin()
	defer tx.Rollback()

	if field == "threads" {
		sqlUpdate := `
			UPDATE Forum SET threads = threads + $1 where lower(slug) like lower($2);
			`
		_, err = tx.Exec(sqlUpdate, number, slug)
		if err != nil {
			return
		}
	} else if field == "posts" {
		sqlUpdate := `
			UPDATE Forum SET posts = posts + $1 where lower(slug) like lower($2);
			`
		_, err = tx.Exec(sqlUpdate, number, slug)
		if err != nil {
			return
		}
	}

	err = tx.Commit()
	if err != nil {
		return
	}
	fmt.Println("database/UpdateFieldsForum +")

	return
}

// isForumUnique checks if there are Players with
func (db DataBase) isForumUnique(slug string, checkUnique *bool) (forum models.Forum, err error) {
	var tx *sql.Tx
	tx, err = db.Db.Begin()
	defer tx.Rollback()
	*checkUnique = false
	sqlStatement := "SELECT f.slug, f.title, f.user " +
		"FROM Forum as f where  lower(f.slug) like lower($1)"

	row := tx.QueryRow(sqlStatement, slug)
	err = row.Scan(&forum.Slug, &forum.Title, &forum.User)
	if err != nil {
		if err == sql.ErrNoRows {
			*checkUnique = true
			err = nil
		}
		fmt.Println("database/isForumUnique Query error")
		return
	}

	err = tx.Commit()
	if err != nil {
		return
	}

	return
}

// CreateForum
func (db *DataBase) CreateForum(forum models.Forum) (forumQuery models.Forum, checkUnique bool, err error) {
	var tx *sql.Tx
	tx, err = db.Db.Begin()
	defer tx.Rollback()

	//forum.User = strings.Replace(forum.User, "_", "\\_", -1)

	checkUnique = false
	if forumQuery, err = db.isForumUnique(forum.Slug, &checkUnique); err != nil {
		fmt.Println("database/CreateUser - fail uniqie")
		return
	}

	if !checkUnique {
		fmt.Println("CreateForum ", forum)
		return
	}

	sqlInsert := `
	INSERT INTO Forum(slug, title, "user") VALUES
    ($1, $2, $3);
		`
	_, err = tx.Exec(sqlInsert, forum.Slug, forum.Title, forum.User)
	if err != nil {
		return
	}
	err = tx.Commit()
	if err != nil {
		return
	}
	//forum.User = strings.Replace(forum.User, "\\_", "_", -1)
	forumQuery = forum

	fmt.Println("database/CreateUser +")

	return
}

// CountForum
func (db *DataBase) CountForum() (count int, err error) {
	var tx *sql.Tx
	tx, err = db.Db.Begin()
	defer tx.Rollback()

	sqlInsert := `
		SELECT COUNT(*) FROM Forum
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
	fmt.Println("database/CountForum +")

	return
}
