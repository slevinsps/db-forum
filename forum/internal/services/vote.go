package api

import (
	"db_forum/internal/models"
	"net/http"
)

func (h *Handler) ThreadVote(rw http.ResponseWriter, r *http.Request) {

	const place = "ForumDetails"

	var (
		err             error
		slugOrID        string
		checkFindThread bool
		checkFindUser   bool
		vote            models.Vote
		thread          models.Thread
	)
	rw.Header().Set("Content-Type", "application/json")

	if slugOrID, err = h.getSlugOrID(r); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		printResult(err, http.StatusBadRequest, place)
		return
	}

	if vote, err = getVoteFromBody(r); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		printResult(err, http.StatusBadRequest, place)
		return
	}

	if thread, checkFindThread, err = h.DB.GetThreadById(slugOrID); err != nil {
		rw.WriteHeader(http.StatusNotFound)
		printResult(err, http.StatusNotFound, place)
		return
	}

	if !checkFindThread {
		rw.WriteHeader(http.StatusNotFound)
		message := models.Message{Message: "Can't find thread by id: " + slugOrID}
		sendJSON(rw, message, place)
		return
	}

	if _, checkFindUser, err = h.DB.GetUserByNickname(vote.Nickname); err != nil {
		rw.WriteHeader(http.StatusNotFound)
		printResult(err, http.StatusNotFound, place)
		return
	}

	if !checkFindUser {
		rw.WriteHeader(http.StatusNotFound)
		message := models.Message{Message: "Can't find user by nickname: " + vote.Nickname}
		sendJSON(rw, message, place)
		return
	}

	if err = h.DB.InsertOrUpdateVoteUser(vote, &thread); err != nil {
		rw.WriteHeader(http.StatusNotFound)
		printResult(err, http.StatusNotFound, place)
		return
	}
	rw.WriteHeader(http.StatusOK)
	sendJSON(rw, thread, place)

	printResult(err, http.StatusCreated, place)
	return
}
