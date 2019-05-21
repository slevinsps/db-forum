package database

import (
	"database/sql"
	"db_forum/internal/models"
	"db_forum/internal/utils"
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
		"FROM Users as u where nickname = $1 or email = $2"

	rows, erro := tx.Query(sqlStatement, nickname, email)
	if erro != nil {
		err = erro
		utils.PrintDebug("database/isNicknameUnique_test Query error")
		return
	}

	defer rows.Close()

	for rows.Next() {
		user := models.User{}
		if err = rows.Scan(&user.About, &user.Email, &user.Fullname,
			&user.Nickname); err != nil {
			utils.PrintDebug("database/isNicknameUnique_test wrong row catched")
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
		"FROM Users as u where email = $1 and u.nickname <> $2"

	rows, erro := tx.Query(sqlStatement, user.Email, user.Nickname)
	if erro != nil {
		err = erro
		utils.PrintDebug("database/isEmailUnique_test Query error")
		return
	}

	defer rows.Close()

	for rows.Next() {
		user := models.User{}
		if err = rows.Scan(&user.About, &user.Email, &user.Fullname,
			&user.Nickname); err != nil {
			utils.PrintDebug("database/isEmailUnique_test wrong row catched")
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
		utils.PrintDebug("database/CreateUser - fail uniqie")
		return
	}

	if len(users) > 0 {
		utils.PrintDebug("CreateUser ", users)
		return
	}

	//utils.PrintDebug(user)
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
	utils.PrintDebug("database/CreateUser +")

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
		utils.PrintDebug("database/CreateUser - fail uniqie")
		return
	}

	if len(users) > 0 {
		utils.PrintDebug("UpdateUser ", users)
		return
	}

	sqlInsert := `
		UPDATE Users SET about = $1, email = $2, fullname = $3, nickname = $4 WHERE nickname = $5
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
	utils.PrintDebug("database/UpdateUser +")

	return
}

// GetUserByNickname
func (db *DataBase) GetUserByNickname(nickname string) (user models.User, checkFindUser bool, err error) {
	checkFindUser = true
	var tx *sql.Tx
	tx, err = db.Db.Begin()
	defer tx.Rollback()

	sqlQuery := `
	SELECT u.about, u.email, u.fullname, u.nickname FROM Users as u where u.nickname = $1;
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

	utils.PrintDebug("database/GetUserByNickname +")

	return
}

// GetUsersByForum
func (db *DataBase) GetUsersByForum(title string, limitStr string, sinceStr string, descStr string) (users []models.User, err error) {

	sqlQuery :=
		"SELECT u.about, u.email, u.fullname, u.nickname " +
			"FROM Users as u join UsersForum as f on f.userNickname = u.nickname and f.forum = $1  "

	if sinceStr != "" {
		if descStr == "true" {
			sqlQuery += "and u.nickname < " + "'" + sinceStr + "'"
		} else {
			sqlQuery += "and " + "'" + sinceStr + "' < u.nickname "
		}
	}

	sqlQuery += " order by 4 "
	if descStr == "true" {
		sqlQuery += "desc "
	}
	if limitStr != "" {
		sqlQuery += "limit " + limitStr + ";"
	}

	utils.PrintDebug(sqlQuery)

	rows, erro := db.Db.Query(sqlQuery, title)

	if erro != nil {
		err = erro
		utils.PrintDebug("database/GetUsersByForum Query error")
		return
	}

	defer rows.Close()

	for rows.Next() {
		user := models.User{}
		if err = rows.Scan(&user.About, &user.Email, &user.Fullname, &user.Nickname); err != nil {
			utils.PrintDebug("database/GetUsersByForum wrong row catched")
			break
		}

		users = append(users, user)
	}
	utils.PrintDebug("database/GetUsersByForum +")

	return
}

// CountUser
func (db *DataBase) CountUser() (count int, err error) {
	var tx *sql.Tx
	tx, err = db.Db.Begin()
	defer tx.Rollback()

	//utils.PrintDebug(user)
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
	utils.PrintDebug("database/CountUser +")

	return
}
