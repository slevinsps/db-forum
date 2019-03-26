package database

import (
	"database/sql"
	"db_forum/internal/models"
	"fmt"

	_ "github.com/lib/pq"
)

// InsertOrUpdateVoteUser
func (db *DataBase) InsertOrUpdateVoteUser(vote models.Vote, thread *models.Thread) (err error) {
	var tx *sql.Tx
	tx, err = db.Db.Begin()
	defer tx.Rollback()

	/*
		sqlInsert := `
		INSERT INTO Vote(nickname, voice) VALUES
		($1, $2) on conflict (nickname) do update set voice = $2;
		`
	*/

	oldVoice := 0
	sqlQuerySelect := `SELECT voice FROM Vote WHERE nickname like $1 and threadId = $2 `
	row := tx.QueryRow(sqlQuerySelect, vote.Nickname, vote.ThreadId)
	err = row.Scan(&oldVoice)
	if err != nil && err != sql.ErrNoRows {
		return
	}

	if err == sql.ErrNoRows {
		sqlQueryInsert := `INSERT INTO Vote(nickname, voice, threadId) VALUES ($1, $2, $3)`
		_, err = tx.Exec(sqlQueryInsert, vote.Nickname, vote.Voice, vote.ThreadId)
		if err != nil {
			return
		}
	} else {
		//fmt.Println(user)
		sqlQueryUpdate := `
		UPDATE Vote SET voice = $1 WHERE nickname like $2 and threadId = $3;
		`
		_, err = tx.Exec(sqlQueryUpdate, vote.Voice, vote.Nickname, vote.ThreadId)
		if err != nil {
			return
		}
	}

	sqlUpdate := `
		UPDATE Thread SET votes = $1 WHERE id = $2 RETURNING votes;
		`
	row2 := tx.QueryRow(sqlUpdate, thread.Votes-oldVoice+vote.Voice, vote.ThreadId)
	err = row2.Scan(&thread.Votes)
	if err != nil {
		return
	}
	err = tx.Commit()
	if err != nil {
		return
	}

	fmt.Println("database/InsertOrUpdateVoteUser +")

	return
}
