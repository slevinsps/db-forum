package database

import (
	"database/sql"
	"db_forum/internal/models"
	"db_forum/internal/utils"
	"time"
)

// GetPostById
func (db *DataBase) GetPostById(id int) (post models.Post, checkFindPost bool, err error) {
	var tx *sql.Tx
	tx, err = db.Db.Begin()
	defer tx.Rollback()
	checkFindPost = true
	sqlQuery :=
		"SELECT p.author, p.created, p.forum, p.is_edited, p.id, p.message, p.parent, p.thread " +
			"FROM Post as p where p.id = $1;"

	row := tx.QueryRow(sqlQuery, id)

	err = row.Scan(&post.Author, &post.Created, &post.Forum, &post.IsEdited, &post.Id, &post.Message, &post.Parent, &post.Thread)
	if err != nil {
		if err == sql.ErrNoRows {
			checkFindPost = false
			err = nil
		}
		utils.PrintDebug("database/GetPostById Scan error")
		return
	}

	err = tx.Commit()
	if err != nil {
		utils.PrintDebug("database/GetPostById Commit error")
		return
	}

	utils.PrintDebug("database/GetPostById +")

	return
}

// CreatePost
func (db *DataBase) CreatePost(posts []models.Post, thread models.Thread, timeNow time.Time) (postQuery []models.Post, check int, err error) {
	var tx *sql.Tx
	tx, err = db.Db.Begin()
	if err != nil {
		return
	}
	defer tx.Rollback()
	var checkPost bool
	var postParent models.Post
	var postQuerySingle models.Post

	check = 0
	if posts[0].Parent != 0 {
		if postParent, checkPost, err = db.GetPostById(posts[0].Parent); err != nil {
			utils.PrintDebug("database/CreatePost - fail checkParent")
			return
		}
		if !checkPost {
			check = -1
			utils.PrintDebug("CreatePost !checkParent")
			return
		}

		if postParent.Thread != thread.Id {
			utils.PrintDebug("CreatePost !postParent")
			check = -1
			return
		}
	}

	sqlInsert := `
	INSERT INTO Post(author, created, forum, message, parent, thread) VALUES ($1, $2, $3, $4, $5, $6) RETURNING author, created, forum, id, message, thread, parent;
	`
	for _, value := range posts {
		value.Thread = thread.Id
		value.Forum = thread.Forum
		row := tx.QueryRow(sqlInsert, value.Author, timeNow, value.Forum, value.Message, value.Parent, value.Thread)
		err = row.Scan(&postQuerySingle.Author, &postQuerySingle.Created, &postQuerySingle.Forum, &postQuerySingle.Id, &postQuerySingle.Message, &postQuerySingle.Thread, &postQuerySingle.Parent)
		if err != nil {
			check = -2
			return
		}
		postQuery = append(postQuery, postQuerySingle)
	}

	err = tx.Commit()
	if err != nil {
		check = -2
		return
	}
	utils.PrintDebug("database/CreatePost +")

	return
}

// GetForumBySlug
func (db *DataBase) GetPostsByThread(thread models.Thread, limitStr string, sinceStr string, sortStr string, descStr string) (posts []models.Post, err error) {
	var tx *sql.Tx
	tx, err = db.Db.Begin()
	defer tx.Rollback()
	var sqlQuery string
	if sortStr == "parent_tree" {
		sqlQuery =
			"SELECT t.author, t.created, t.forum, t.id, t.message, t.thread, t.parent " +
				"FROM Post as t where branch in (SELECT p.id FROM Post as p WHERE p.thread = $1 AND parent = 0 "
		if sinceStr != "" {
			if descStr == "true" {
				sqlQuery += " and p.id < (SELECT branch FROM Post WHERE id = " + sinceStr + ") "
			} else {
				sqlQuery += " and p.id > (SELECT branch FROM Post WHERE id =  " + sinceStr + ") "
			}

		}
		sqlQuery += " order by p.id "
		if descStr == "true" {
			sqlQuery += "desc "
		}
		if limitStr != "" {
			sqlQuery += "limit " + limitStr
		}
		if descStr == "true" {
			sqlQuery += ") order by t.branch desc, t.path "
		} else {
			sqlQuery += ") order by t.path; "
		}

	} else if sortStr == "tree" {
		sqlQuery =
			"SELECT t.author, t.created, t.forum, t.id, t.message, t.thread, t.parent " +
				"FROM Post as t where t.thread = $1 "
		if sinceStr != "" {
			if descStr == "true" {
				sqlQuery += " and t.path < (SELECT path FROM Post WHERE id = " + sinceStr + ") "
			} else {
				sqlQuery += " and t.path > (SELECT path FROM Post WHERE id =  " + sinceStr + ") "
			}

		}

		sqlQuery += " order by t.path "
		if descStr == "true" {
			sqlQuery += "desc "
		}

		if limitStr != "" {
			sqlQuery += "limit " + limitStr + ";"
		}
	} else {
		sqlQuery =
			"SELECT t.author, t.created, t.forum, t.id, t.message, t.thread, t.parent " +
				"FROM Post as t where t.thread = $1 "
		if sinceStr != "" {
			if descStr == "true" {
				sqlQuery += " and t.id < " + sinceStr + " "
			} else {
				sqlQuery += " and t.id > " + sinceStr + " "
			}

		}
		sqlQuery += " order by t.id "
		if descStr == "true" {
			sqlQuery += "desc "
		}

		if limitStr != "" {
			sqlQuery += "limit " + limitStr + ";"
		}
	}
	rows, erro := tx.Query(sqlQuery, thread.Id)
	if erro != nil {
		err = erro
		utils.PrintDebug("database/GetPostsByThread Query error")
		return
	}

	defer rows.Close()

	for rows.Next() {
		post := models.Post{}
		if err = rows.Scan(&post.Author, &post.Created, &post.Forum, &post.Id, &post.Message, &post.Thread, &post.Parent); err != nil {
			utils.PrintDebug("database/GetPostsByThread wrong row catched")
			break
		}

		posts = append(posts, post)
	}

	err = tx.Commit()
	if err != nil {
		return
	}

	utils.PrintDebug("database/GetPostsByThread +")

	return
}

// UpdateThread
func (db *DataBase) UpdatePost(message string, isEdit bool, id int) (post models.Post, err error) {
	var tx *sql.Tx
	tx, err = db.Db.Begin()
	defer tx.Rollback()
	var sqlQuery string
	if isEdit {
		sqlQuery = `
		UPDATE Post SET message = $1, is_edited = true where id = $2 RETURNING author, created, forum, id, is_edited, message, parent, thread;
		`
	} else {
		sqlQuery = `
		UPDATE Post SET message = $1 where id = $2 RETURNING author, created, forum, id, is_edited, message, parent, thread;
		`
	}
	row := tx.QueryRow(sqlQuery, message, id)
	err = row.Scan(&post.Author, &post.Created, &post.Forum, &post.Id, &post.IsEdited, &post.Message, &post.Parent, &post.Thread)

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

// CountPost
func (db *DataBase) CountPost() (count int, err error) {
	var tx *sql.Tx
	tx, err = db.Db.Begin()
	defer tx.Rollback()

	sqlInsert := `
		SELECT COUNT(*) FROM Post
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
	utils.PrintDebug("database/CountPost +")

	return
}
