package database

import (
	"database/sql"
	"db_forum/internal/models"
	"fmt"

	_ "github.com/lib/pq"
)

// InsertVoteUser
func (db *DataBase) InsertVoteUser(nickname string, voice int, threadID int, outErr chan error) {
	var (
		tx  *sql.Tx
		err error
	)
	tx, err = db.Db.Begin()
	if err != nil {
		outErr <- err
		return
	}
	defer tx.Rollback()
	sqlQueryInsert := `INSERT INTO Vote(nickname, voice, threadId) VALUES ($1, $2, $3);`
	_, err = tx.Exec(sqlQueryInsert, nickname, voice, threadID)
	if err != nil {
		outErr <- err
		return
	}
	err = tx.Commit()
	outErr <- err
}

// UpdateVoteUser
func (db *DataBase) UpdateVoteUser(nickname string, voice int, threadID int, outErr chan error) {
	var (
		tx  *sql.Tx
		err error
	)
	tx, err = db.Db.Begin()
	if err != nil {
		outErr <- err
		return
	}
	defer tx.Rollback()
	sqlQueryUpdate := `UPDATE Vote SET voice = $1 WHERE nickname like $2 and threadId = $3;`
	_, err = tx.Exec(sqlQueryUpdate, voice, nickname, threadID)
	if err != nil {
		outErr <- err
		return
	}
	err = tx.Commit()
	outErr <- err
}

// InsertOrUpdateVoteUser
func (db *DataBase) InsertOrUpdateVoteUser(vote models.Vote, thread *models.Thread) (err error) {
	var tx *sql.Tx
	tx, err = db.Db.Begin()
	defer tx.Rollback()

	oldVoice := 0
	sqlQuerySelect := `SELECT voice FROM Vote WHERE nickname like $1 and threadId = $2 `
	row := tx.QueryRow(sqlQuerySelect, vote.Nickname, vote.ThreadId)
	err = row.Scan(&oldVoice)
	if err != nil && err != sql.ErrNoRows {
		return
	}

	outErr := make(chan error, 1)
	defer close(outErr)

	if err == sql.ErrNoRows {
		go db.InsertVoteUser(vote.Nickname, vote.Voice, vote.ThreadId, outErr)
	} else {
		go db.UpdateVoteUser(vote.Nickname, vote.Voice, vote.ThreadId, outErr)
	}

	sqlUpdate := `
		UPDATE Thread SET votes = $1 WHERE id = $2;
		`
	_, err = tx.Exec(sqlUpdate, thread.Votes-oldVoice+vote.Voice, vote.ThreadId)
	if err != nil {
		return
	}

	thread.Votes = thread.Votes - oldVoice + vote.Voice
	err = tx.Commit()
	if err != nil {
		return
	}

	err = <-outErr
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("database/InsertOrUpdateVoteUser +")

	return
}
