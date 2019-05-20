package api

import (
	"encoding/json"
	"db_forum/internal/models"
	"db_forum/internal/utils"
	"net/http"
	"sync"
)

func (h *Handler) ThreadCreate(rw http.ResponseWriter, r *http.Request) {
	const place = "ThreadCreate"

	utils.PrintDebug("ThreadCreate---------------------------begin")
	var (
		err            error
		thread         models.Thread
		checkUnique    bool
		slug           string
		forum          models.Forum
		checkFindForum bool
	)

	if slug, err = h.getSlug(r); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		printResult(err, http.StatusBadRequest, place)
		return
	}

	if thread, err = getThreadFromBody(r); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		printResult(err, http.StatusBadRequest, place)
		return
	}
	utils.PrintDebug("thread---------------------------end")

	rw.Header().Set("Content-Type", "application/json")

	if forum, checkFindForum, err = h.DB.GetForumBySlug(slug); err != nil {
		rw.WriteHeader(http.StatusNotFound)
		printResult(err, http.StatusNotFound, place)
		return
	}
	if !checkFindForum {
		rw.WriteHeader(http.StatusNotFound)
		message := models.Message{Message: "Can't find forum by slag: " + slug}
		resBytes, _ := message.MarshalJSON()
		sendJSON(rw, resBytes, place)
		return
	}

	thread.Forum = forum.Slug

	if thread, checkUnique, err = h.DB.CreateThread(thread); err != nil {
		rw.WriteHeader(http.StatusNotFound)
		message := models.Message{Message: "Can't find user"}
		resBytes, _ := message.MarshalJSON()
		sendJSON(rw, resBytes, place)
		printResult(err, http.StatusNotFound, place)
		return
	}

	if checkUnique {
		rw.WriteHeader(http.StatusCreated)
		resBytes, _ := thread.MarshalJSON()
		sendJSON(rw, resBytes, place)
	} else {
		rw.WriteHeader(http.StatusConflict)
		resBytes, _ := thread.MarshalJSON()
		sendJSON(rw, resBytes, place)
	}

	utils.PrintDebug("ThreadCreate---------------------------end")
	printResult(err, http.StatusCreated, place)
	return
}

func getUserByNickname(wg *sync.WaitGroup, nickname string, h *Handler, out chan models.User, outErr chan int) {
	defer wg.Done()
	var (
		err           error
		user          models.User
		checkFindUser bool
	)

	if user, checkFindUser, err = h.DB.GetUserByNickname(nickname); err != nil {
		outErr <- -1
		return
	}

	if checkFindUser {
		out <- user
		outErr <- 0
	} else {
		outErr <- -2
	}
}

func (h *Handler) ForumThreads(rw http.ResponseWriter, r *http.Request) {

	const place = "ForumThreads"

	var (
		err            error
		slug           string
		forum          models.Forum
		threads        []models.Thread
		checkFindForum bool
	)

	if slug, err = h.getSlug(r); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		printResult(err, http.StatusBadRequest, place)
		return
	}

	if forum, checkFindForum, err = h.DB.GetForumBySlug(slug); err != nil {
		rw.WriteHeader(http.StatusNotFound)
		printResult(err, http.StatusNotFound, place)
		return
	}

	rw.Header().Set("Content-Type", "application/json")

	if checkFindForum {
		if threads, err = h.DB.GetThreadsByForum(forum.Slug, r.FormValue("limit"), r.FormValue("since"), r.FormValue("desc")); err != nil {
			rw.WriteHeader(http.StatusNotFound)
			printResult(err, http.StatusNotFound, place)
			return
		}
		rw.WriteHeader(http.StatusOK)
		if len(threads) == 0 {
			rw.Write([]byte("[]"))
		} else {

			resBytes, _ := json.Marshal(threads)
			sendJSON(rw, resBytes, place)
		}
	} else {
		rw.WriteHeader(http.StatusNotFound)
		message := models.Message{Message: "Can't find forum by sluq: " + slug}
		resBytes, _ := message.MarshalJSON()
		sendJSON(rw, resBytes, place)
	}

	printResult(err, http.StatusCreated, place)
	return
}

func (h *Handler) ThreadDetails(rw http.ResponseWriter, r *http.Request) {

	const place = "ThreadDetails"

	var (
		err             error
		slugOrID        string
		checkFindThread bool
		thread          models.Thread
	)

	if slugOrID, err = h.getSlugOrID(r); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		printResult(err, http.StatusBadRequest, place)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	if thread, checkFindThread, err = h.DB.GetThreadById(slugOrID); err != nil {
		rw.WriteHeader(http.StatusNotFound)
		printResult(err, http.StatusNotFound, place)
		return
	}

	if !checkFindThread {
		rw.WriteHeader(http.StatusNotFound)
		message := models.Message{Message: "Can't find thread by id: " + slugOrID}
		resBytes, _ := message.MarshalJSON()
		sendJSON(rw, resBytes, place)
		return
	} else {
		rw.WriteHeader(http.StatusOK)
		resBytes, _ := thread.MarshalJSON()
		sendJSON(rw, resBytes, place)
	}

	printResult(err, http.StatusCreated, place)
	return
}

func (h *Handler) ThreadUpdate(rw http.ResponseWriter, r *http.Request) {

	const place = "ThreadUpdate"

	var (
		err             error
		slugOrID        string
		checkFindThread bool
		threadNew       models.Thread
		threadOld       models.Thread
	)
	rw.Header().Set("Content-Type", "application/json")

	if slugOrID, err = h.getSlugOrID(r); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		printResult(err, http.StatusBadRequest, place)
		return
	}

	if threadNew, err = getThreadFromBody(r); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		printResult(err, http.StatusBadRequest, place)
		return
	}

	if threadOld, checkFindThread, err = h.DB.GetThreadById(slugOrID); err != nil {
		rw.WriteHeader(http.StatusNotFound)
		printResult(err, http.StatusNotFound, place)
		return
	}

	if !checkFindThread {
		rw.WriteHeader(http.StatusNotFound)
		message := models.Message{Message: "Can't find thread by id: " + slugOrID}
	
		resBytes, _ := message.MarshalJSON()
		sendJSON(rw, resBytes, place)
		
		return
	}

	if threadOld, err = h.DB.UpdateThread(threadNew, threadOld); err != nil {
		rw.WriteHeader(http.StatusNotFound)
		printResult(err, http.StatusNotFound, place)
		return
	}
	rw.WriteHeader(http.StatusOK)
	resBytes, _ := threadOld.MarshalJSON()
	sendJSON(rw, resBytes, place)
	printResult(err, http.StatusCreated, place)
	return
}
