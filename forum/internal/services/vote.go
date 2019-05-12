package api

import (
	"sync"
	"db_forum/internal/models"
	"net/http"
)

func (h *Handler) ThreadVote(rw http.ResponseWriter, r *http.Request) {

	const place = "ForumDetails"

	var (
		err             error
		slugOrID        string
		vote            models.Vote
		thread          models.Thread
		user            models.User
	)
	rw.Header().Set("Content-Type", "application/json")
	wg := &sync.WaitGroup{}
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

	threadChan := make(chan models.Thread, 1)
	threadChanErr := make(chan int, 1)
	userChan := make(chan models.User, 1)
	userChanErr := make(chan int, 1)
	defer close(threadChan)
	defer close(threadChanErr)
	defer close(userChan)
	defer close(userChanErr)

	wg.Add(2)
	go getThreadByIdThread(wg, slugOrID, h, threadChan, threadChanErr)
	go getUserByNickname(wg, vote.Nickname, h, userChan, userChanErr)
	

	wg.Wait()
	errUser := <-userChanErr
	errTread := <-threadChanErr

	if (errUser == -1 || errTread == -1) {
		rw.WriteHeader(http.StatusNotFound)
		printResult(err, http.StatusNotFound, place)
		return
	}

	if (errUser == -2) {
		rw.WriteHeader(http.StatusNotFound)
		message := models.Message{Message: "Can't find user by nickname: " + vote.Nickname}
		sendJSON(rw, message, place)
		return
	} else {
		user = <-userChan
	}

	if (errTread == -2) {
		rw.WriteHeader(http.StatusNotFound)
		message := models.Message{Message: "Can't find thread by id: " + slugOrID}
		sendJSON(rw, message, place)
		return
	} else {
		thread = <-threadChan
	}
	vote.ThreadId = thread.Id
	vote.Nickname = user.Nickname
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


func getThreadByIdThread(wg *sync.WaitGroup, slugOrID string, h *Handler, out chan models.Thread, outErr chan int) {
	defer wg.Done()
	var (
		err error
		thread models.Thread
		checkFindThread bool
	)

	if thread, checkFindThread, err = h.DB.GetThreadById(slugOrID); err != nil {
		outErr <- -1;
		return
	}
	if (checkFindThread) {
		out <- thread
		outErr <- 0;
	} else {
		outErr <- -2;
	}
}

func getUserByNickname(wg *sync.WaitGroup, nickname string, h *Handler, out chan models.User, outErr chan int) {
	defer wg.Done()
	var (
		err error
		user models.User
		checkFindUser bool
	)

	if user, checkFindUser, err = h.DB.GetUserByNickname(nickname); err != nil {
		outErr <- -1;
		return
	}

	if (checkFindUser) {
		out <- user
		outErr <- 0;
	} else {
		outErr <- -2;
	}
}