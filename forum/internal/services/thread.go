package api

import (
	"db_forum/internal/models"
	"fmt"
	"net/http"
	"sync"
)

func (h *Handler) ThreadCreate(rw http.ResponseWriter, r *http.Request) {
	const place = "ThreadCreate"

	fmt.Println("ThreadCreate---------------------------begin")
	var (
		err         error
		thread      models.Thread
		checkUnique bool
		slug        string
		forum       models.Forum
	)

	wg := &sync.WaitGroup{}
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
	fmt.Println("thread---------------------------end")

	rw.Header().Set("Content-Type", "application/json")

	forumChan := make(chan models.Forum, 1)
	forumChanErr := make(chan int, 1)
	userChan := make(chan models.User, 1)
	userChanErr := make(chan int, 1)
	defer close(forumChan)
	defer close(forumChanErr)
	defer close(userChan)
	defer close(userChanErr)

	wg.Add(2)
	go getForumBySlug(wg, slug, h, forumChan, forumChanErr)
	go getUserByNickname(wg, thread.Author, h, userChan, userChanErr)

	wg.Wait()
	errUser := <-userChanErr
	errForum := <-forumChanErr

	if errUser == -1 || errForum == -1 {
		rw.WriteHeader(http.StatusNotFound)
		printResult(err, http.StatusNotFound, place)
		return
	}

	if errUser == -2 {
		rw.WriteHeader(http.StatusNotFound)
		message := models.Message{Message: "Can't find user by nickname: " + thread.Author}
		sendJSON(rw, message, place)
		return
	} else {
		<-userChan
	}

	if errForum == -2 {
		rw.WriteHeader(http.StatusNotFound)
		message := models.Message{Message: "Can't find forum by slag: " + slug}
		sendJSON(rw, message, place)
		return
	} else {
		forum = <-forumChan
	}

	thread.Forum = forum.Slug

	if thread, checkUnique, err = h.DB.CreateThread(thread); err != nil {
		rw.WriteHeader(http.StatusNotFound)
		printResult(err, http.StatusNotFound, place)
		return
	}

	if checkUnique {
		rw.WriteHeader(http.StatusCreated)
		sendJSON(rw, thread, place)
	} else {
		rw.WriteHeader(http.StatusConflict)
		sendJSON(rw, thread, place)
	}

	fmt.Println("ThreadCreate---------------------------end")
	printResult(err, http.StatusCreated, place)
	return
}

func getForumBySlug(wg *sync.WaitGroup, slug string, h *Handler, out chan models.Forum, outErr chan int) {
	defer wg.Done()
	var (
		err            error
		forum          models.Forum
		checkFindForum bool
	)

	if forum, checkFindForum, err = h.DB.GetForumBySlug(slug); err != nil {
		outErr <- -1
		return
	}
	if checkFindForum {
		out <- forum
		outErr <- 0
	} else {
		outErr <- -2
	}
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
