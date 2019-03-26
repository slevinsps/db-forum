package database

import (
	"database/sql"
	"db_forum/internal/models"
	"fmt"
)

type DataBase struct {
	Db *sql.DB
}

// isNicknameEmailUnique_test checks if there are Players with
func (db DataBase) isNicknameEmailUnique_test(nickname string, email string, users *[]models.User) (err error) {
	var tx *sql.Tx
	tx, err = db.Db.Begin()
	defer tx.Rollback()
	sqlStatement := "SELECT u.about, u.email, u.fullname, u.nickname " +
		"FROM Users as u where nickname like $1 or lower(email) like lower($2)"

	rows, erro := tx.Query(sqlStatement, nickname, email)
	if erro != nil {
		err = erro
		fmt.Println("database/isNicknameUnique_test Query error")
		return
	}

	defer rows.Close()

	for rows.Next() {
		user := models.User{}
		if err = rows.Scan(&user.About, &user.Email, &user.Fullname,
			&user.Nickname); err != nil {
			fmt.Println("database/isNicknameUnique_test wrong row catched")
			break
		}

		*users = append(*users, user)
	}

	err = tx.Commit()
	if err != nil {
		return
	}

	return
}

func (db DataBase) isEmailUnique_test(user models.User, users *[]models.User) (err error) {
	var tx *sql.Tx
	tx, err = db.Db.Begin()
	defer tx.Rollback()
	sqlStatement := "SELECT u.about, u.email, u.fullname, u.nickname " +
		"FROM Users as u where lower(email) like lower($1) and u.nickname not like $2"

	rows, erro := tx.Query(sqlStatement, user.Email, user.Nickname)
	if erro != nil {
		err = erro
		fmt.Println("database/isEmailUnique_test Query error")
		return
	}

	defer rows.Close()

	for rows.Next() {
		user := models.User{}
		if err = rows.Scan(&user.About, &user.Email, &user.Fullname,
			&user.Nickname); err != nil {
			fmt.Println("database/isEmailUnique_test wrong row catched")
			break
		}

		*users = append(*users, user)
	}

	err = tx.Commit()
	if err != nil {
		return
	}

	return
}

// CreateUser
func (db *DataBase) CreateUser(user models.User) (users []models.User, checkUnique bool, err error) {

	var tx *sql.Tx
	tx, err = db.Db.Begin()
	defer tx.Rollback()

	checkUnique = false
	if err = db.isNicknameEmailUnique_test(user.Nickname, user.Email, &users); err != nil {
		fmt.Println("database/CreateUser - fail uniqie")
		return
	}

	if len(users) > 0 {
		fmt.Println("CreateUser ", users)
		return
	}

	//fmt.Println(user)
	sqlInsert := `
	INSERT INTO Users(about, email, fullname, nickname) VALUES
    ($1, $2, $3, $4);
		`
	_, err = tx.Exec(sqlInsert, user.About, user.Email, user.Fullname, user.Nickname)
	if err != nil {
		return
	}
	err = tx.Commit()
	if err != nil {
		return
	}
	users = append(users, user)
	checkUnique = true
	fmt.Println("database/CreateUser +")

	return
}

// UpdateUser
func (db *DataBase) UpdateUser(user models.User) (checkUnique bool, err error) {

	var tx *sql.Tx
	var users []models.User
	tx, err = db.Db.Begin()
	defer tx.Rollback()
	checkUnique = false
	if err = db.isEmailUnique_test(user, &users); err != nil {
		fmt.Println("database/CreateUser - fail uniqie")
		return
	}

	if len(users) > 0 {
		fmt.Println("UpdateUser ", users)
		return
	}

	sqlInsert := `
		UPDATE Users SET about = $1, email = $2, fullname = $3, nickname = $4 WHERE nickname like $5
		`
	_, err = tx.Exec(sqlInsert, user.About, user.Email, user.Fullname, user.Nickname, user.Nickname)
	if err != nil {
		return
	}
	err = tx.Commit()
	if err != nil {
		return
	}
	checkUnique = true
	fmt.Println("database/UpdateUser +")

	return
}

// GetUserByNickname
func (db *DataBase) GetUserByNickname(nickname string) (user models.User, checkFindUser bool, err error) {
	checkFindUser = true
	var tx *sql.Tx
	tx, err = db.Db.Begin()
	defer tx.Rollback()

	sqlQuery := `
	SELECT u.about, u.email, u.fullname, u.nickname FROM Users as u where u.nickname like $1;
	`

	row := tx.QueryRow(sqlQuery, nickname)
	err = row.Scan(&user.About, &user.Email, &user.Fullname, &user.Nickname)
	if err != nil {
		if err == sql.ErrNoRows {
			checkFindUser = false
			err = nil
		}
		return
	}

	err = tx.Commit()
	if err != nil {
		return
	}

	fmt.Println("database/GetUserByNickname +")

	return
}

// GetUsersByForum
func (db *DataBase) GetUsersByForum(title string, limitStr string, sinceStr string, descStr string) (users []models.User, err error) {
	var tx *sql.Tx
	tx, err = db.Db.Begin()
	defer tx.Rollback()

	sqlQuery :=
		"SELECT distinct u.about, u.email, u.fullname, u.nickname " +
			"FROM Users as u join Post as p on p.author like u.nickname where lower(p.forum) like lower($1)  "

	if sinceStr != "" {
		if descStr == "true" {
			sqlQuery += "and u.nickname < " + "'" + sinceStr + "'"
		} else {
			sqlQuery += "and u.nickname > " + "'" + sinceStr + "'"
		}
	}

	sqlQuery += " union SELECT distinct u.about, u.email, u.fullname, u.nickname " +
		"FROM Users as u join Thread as t on t.author like u.nickname where lower(t.forum) like lower($1) "

	if sinceStr != "" {
		if descStr == "true" {
			sqlQuery += "and u.nickname < " + "'" + sinceStr + "'"
		} else {
			sqlQuery += "and u.nickname > " + "'" + sinceStr + "'"
		}
	}

	sqlQuery += " order by 4 "
	if descStr == "true" {
		sqlQuery += "desc "
	}
	if limitStr != "" {
		sqlQuery += "limit " + limitStr + ";"
	}

	fmt.Println(sqlQuery)

	rows, erro := tx.Query(sqlQuery, title)
	if erro != nil {
		err = erro
		fmt.Println("database/GetUsersByForum Query error")
		return
	}

	defer rows.Close()

	for rows.Next() {
		user := models.User{}
		if err = rows.Scan(&user.About, &user.Email, &user.Fullname, &user.Nickname); err != nil {
			fmt.Println("database/GetUsersByForum wrong row catched")
			break
		}

		users = append(users, user)
	}

	err = tx.Commit()
	if err != nil {
		return
	}

	fmt.Println("database/GetUsersByForum +")

	return
}

// CountUser
func (db *DataBase) CountUser() (count int, err error) {
	var tx *sql.Tx
	tx, err = db.Db.Begin()
	defer tx.Rollback()

	//fmt.Println(user)
	sqlInsert := `
		SELECT COUNT(*) FROM Users
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
	fmt.Println("database/CountUser +")

	return
}
