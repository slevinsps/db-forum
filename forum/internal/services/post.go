package api

import (
	"db_forum/internal/models"
	"db_forum/internal/utils"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func (h *Handler) ThreadCreatePost(rw http.ResponseWriter, r *http.Request) {
	const place = "ThreadCreatePost"

	var (
		err             error
		posts           []models.Post
		checkFindThread bool
		thread          models.Thread
		slugOrID        string
		check           int
	)
	rw.Header().Set("Content-Type", "application/json")

	if slugOrID, err = h.getSlugOrID(r); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		printResult(err, http.StatusBadRequest, place)
		return
	}

	if posts, err = getPostsFromBody(r); err != nil {
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

	if len(posts) == 0 {
		rw.WriteHeader(http.StatusCreated)
		rw.Write([]byte("[]"))
		return
	}

	timeNow := time.Now()

	if posts, check, err = h.DB.CreatePost(posts, thread, timeNow); err != nil {
		rw.WriteHeader(http.StatusNotFound)
		message := models.Message{Message: "Can't find post author"}
		sendJSON(rw, message, place)
		printResult(err, http.StatusNotFound, place)
		return
	}
	if check == -1 {
		rw.WriteHeader(http.StatusConflict)
		message := models.Message{Message: "Parent post was created in another thread"}
		sendJSON(rw, message, place)
		return
	} else if check == -2 {
		rw.WriteHeader(http.StatusNotFound)
		message := models.Message{Message: "Can't find post author"}
		sendJSON(rw, message, place)
		return
	}

	_ = h.DB.UpdateFieldsForum(thread.Forum, len(posts), "posts")
	rw.WriteHeader(http.StatusCreated)
	sendJSON(rw, posts, place)
	printResult(err, http.StatusCreated, place)

	return
}

func (h *Handler) ThreadPosts(rw http.ResponseWriter, r *http.Request) {

	const place = "ThreadPosts"

	var (
		err             error
		slugOrID        string
		thread          models.Thread
		checkFindThread bool
		posts           []models.Post
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
	}

	if posts, err = h.DB.GetPostsByThread(thread, r.FormValue("limit"), r.FormValue("since"), r.FormValue("sort"), r.FormValue("desc")); err != nil {
		rw.WriteHeader(http.StatusNotFound)
		printResult(err, http.StatusNotFound, place)
		return
	}

	rw.WriteHeader(http.StatusOK)
	if len(posts) == 0 {
		rw.Write([]byte("[]"))
	} else {
		sendJSON(rw, posts, place)
	}

	printResult(err, http.StatusCreated, place)
	return
}

func (h *Handler) PostDetails(rw http.ResponseWriter, r *http.Request) {

	const place = "PostDetails"

	var (
		err           error
		idStr         string
		id            int
		checkFindPost bool
		post          models.Post
		user          models.User
		forum         models.Forum
		thread        models.Thread
		details       models.Details
	)
	//r.ParseForm()
	related := strings.Split(r.FormValue("related"), ",")
	if idStr, err = h.getId(r); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		printResult(err, http.StatusBadRequest, place)
		return
	}
	id, err = strconv.Atoi(idStr)
	if err != nil {
		utils.PrintDebug("strconv error in PostDetails")
		return
	}
	rw.Header().Set("Content-Type", "application/json")
	if post, checkFindPost, err = h.DB.GetPostById(id); err != nil {
		rw.WriteHeader(http.StatusNotFound)
		printResult(err, http.StatusNotFound, place)
		return
	}
	if !checkFindPost {
		rw.WriteHeader(http.StatusNotFound)
		message := models.Message{Message: "Can't find post by id: " + strconv.Itoa(id)}
		sendJSON(rw, message, place)
		return
	} else {

		details.Post = &post
	}

	for _, item := range related {
		if item == "user" {
			if user, _, err = h.DB.GetUserByNickname(post.Author); err != nil {
				rw.WriteHeader(http.StatusNotFound)
				printResult(err, http.StatusNotFound, place)
				return
			}
			details.Author = &user
		}
		if item == "thread" {
			if thread, _, err = h.DB.GetThreadById(strconv.Itoa(post.Thread)); err != nil {
				rw.WriteHeader(http.StatusNotFound)
				printResult(err, http.StatusNotFound, place)
				return
			}
			details.Thread = &thread
		}
		if item == "forum" {
			if forum, _, err = h.DB.GetForumBySlug(post.Forum); err != nil {
				rw.WriteHeader(http.StatusNotFound)
				printResult(err, http.StatusNotFound, place)
				return
			}
			details.Forum = &forum
		}

	}
	rw.WriteHeader(http.StatusOK)
	sendJSON(rw, details, place)

	printResult(err, http.StatusCreated, place)
	return
}

func (h *Handler) PostUpdate(rw http.ResponseWriter, r *http.Request) {

	const place = "PostUpdate"

	var (
		err           error
		idStr         string
		id            int
		post          models.Post
		message       models.Message
		checkFindPost bool
	)
	rw.Header().Set("Content-Type", "application/json")

	if idStr, err = h.getId(r); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		printResult(err, http.StatusBadRequest, place)
		return
	}
	id, err = strconv.Atoi(idStr)
	if err != nil {
		utils.PrintDebug("strconv error in PostDetails")
		return
	}

	if message, err = getMessageFromBody(r); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		printResult(err, http.StatusBadRequest, place)
		return
	}

	if post, checkFindPost, err = h.DB.GetPostById(id); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		printResult(err, http.StatusBadRequest, place)
		return
	}

	if !checkFindPost {
		rw.WriteHeader(http.StatusNotFound)
		message := models.Message{Message: "Can't find post by id: " + strconv.Itoa(id)}
		sendJSON(rw, message, place)
		return
	}

	if message.Message == "" {
		rw.WriteHeader(http.StatusOK)
		sendJSON(rw, post, place)
		printResult(err, http.StatusCreated, place)
		return
	}
	if post, err = h.DB.UpdatePost(message.Message, message.Message != post.Message, id); err != nil {
		rw.WriteHeader(http.StatusNotFound)
		printResult(err, http.StatusNotFound, place)
		return
	}
	rw.WriteHeader(http.StatusOK)
	sendJSON(rw, post, place)

	printResult(err, http.StatusCreated, place)
	return
}
