package database

import (
	"database/sql"
	"db_forum/internal/models"
	"db_forum/internal/utils"

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
	sqlQueryUpdate := `UPDATE Vote SET voice = $1 WHERE nickname = $2 and threadId = $3;`
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

	sqlQuerySelect := `INSERT INTO Vote(nickname, threadId, voice) VALUES ($1, $2, $3) ` +
		`on conflict (nickname, threadId) do update set voice = $3 RETURNING voicePrevious;`
	row := tx.QueryRow(sqlQuerySelect, vote.Nickname, vote.ThreadId, vote.Voice)
	err = row.Scan(&oldVoice)
	if err != nil && err != sql.ErrNoRows {
		return
	}
	newVoice := vote.Voice
	thread.Votes = thread.Votes - oldVoice + newVoice
	err = tx.Commit()
	if err != nil {
		return
	}

	utils.PrintDebug("database/InsertOrUpdateVoteUser +")

	return
}
