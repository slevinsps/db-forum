package api

import (
	"db_forum/internal/models"
	"net/http"
)

func (h *Handler) ThreadCreate(rw http.ResponseWriter, r *http.Request) {
	const place = "ThreadCreate"

	var (
		err            error
		thread         models.Thread
		checkUnique    bool
		checkFindUser  bool
		checkFindForum bool
		slug           string
		forum          models.Forum
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

	if _, checkFindUser, err = h.DB.GetUserByNickname(thread.Author); err != nil {
		rw.WriteHeader(http.StatusNotFound)
		printResult(err, http.StatusNotFound, place)
		return
	}

	rw.Header().Set("Content-Type", "application/json")

	if !checkFindUser {
		rw.WriteHeader(http.StatusNotFound)
		message := models.Message{Message: "Can't find user by nickname: " + thread.Author}
		sendJSON(rw, message, place)
		return
	}

	if forum, checkFindForum, err = h.DB.GetForumBySlug(slug); err != nil {
		rw.WriteHeader(http.StatusNotFound)
		printResult(err, http.StatusNotFound, place)
		return
	}

	if !checkFindForum {
		rw.WriteHeader(http.StatusNotFound)
		message := models.Message{Message: "Can't find forum by slag: " + slug}
		sendJSON(rw, message, place)
		return
	}
	thread.Forum = forum.Slug

	if thread, checkUnique, err = h.DB.CreateThread(thread, slug); err != nil {
		rw.WriteHeader(http.StatusNotFound)
		printResult(err, http.StatusNotFound, place)
		return
	}
	rw.Header().Set("Content-Type", "application/json")

	if checkUnique {
		rw.WriteHeader(http.StatusCreated)
		sendJSON(rw, thread, place)
	} else {
		rw.WriteHeader(http.StatusConflict)
		sendJSON(rw, thread, place)
	}

	printResult(err, http.StatusCreated, place)
	return
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
			sendJSON(rw, threads, place)
		}
	} else {
		rw.WriteHeader(http.StatusNotFound)
		message := models.Message{Message: "Can't find forum by sluq: " + slug}
		sendJSON(rw, message, place)
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
		sendJSON(rw, message, place)
		return
	} else {
		rw.WriteHeader(http.StatusOK)
		sendJSON(rw, thread, place)
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
		sendJSON(rw, message, place)
		return
	}

	if threadOld, err = h.DB.UpdateThread(threadNew, threadOld); err != nil {
		rw.WriteHeader(http.StatusNotFound)
		printResult(err, http.StatusNotFound, place)
		return
	}
	rw.WriteHeader(http.StatusOK)
	sendJSON(rw, threadOld, place)

	printResult(err, http.StatusCreated, place)
	return
}
