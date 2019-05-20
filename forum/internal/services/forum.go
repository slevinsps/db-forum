package api

import (
	"db_forum/internal/database"
	"db_forum/internal/models"
	"net/http"
)

type Handler struct {
	DB database.DataBase
}

func (h *Handler) ForumCreate(rw http.ResponseWriter, r *http.Request) {
	const place = "ForumCreate"

	var (
		err           error
		forum         models.Forum
		checkUnique   bool
		checkFindUser bool
		user          models.User
	)

	if forum, err = getForumFromBody(r); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		printResult(err, http.StatusBadRequest, place)
		return
	}

	if user, checkFindUser, err = h.DB.GetUserByNickname(forum.User); err != nil {
		rw.WriteHeader(http.StatusNotFound)
		printResult(err, http.StatusNotFound, place)
		return
	}

	forum.User = user.Nickname

	rw.Header().Set("Content-Type", "application/json")

	if !checkFindUser {
		rw.WriteHeader(http.StatusNotFound)
		message := models.Message{Message: "Can't find user by nickname: " + forum.User}
		resBytes, _ := message.MarshalJSON()
		sendJSON(rw, resBytes, place)
		return
	}

	if forum, checkUnique, err = h.DB.CreateForum(forum); err != nil {
		rw.WriteHeader(http.StatusNotFound)
		printResult(err, http.StatusNotFound, place)
		return
	}
	rw.Header().Set("Content-Type", "application/json")

	if checkUnique {
		rw.WriteHeader(http.StatusCreated)
		resBytes, _ := forum.MarshalJSON()
		sendJSON(rw, resBytes, place)
	} else {
		rw.WriteHeader(http.StatusConflict)
		resBytes, _ := forum.MarshalJSON()
		sendJSON(rw, resBytes, place)
	}

	printResult(err, http.StatusCreated, place)
	return
}

func (h *Handler) ForumDetails(rw http.ResponseWriter, r *http.Request) {

	const place = "ForumDetails"

	var (
		err            error
		slug           string
		forum          models.Forum
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
		rw.WriteHeader(http.StatusOK)
		resBytes, _ := forum.MarshalJSON()
		sendJSON(rw, resBytes, place)
	} else {
		rw.WriteHeader(http.StatusNotFound)
		message := models.Message{Message: "Can't find forum by sluq: " + slug}
		resBytes, _ := message.MarshalJSON()
		sendJSON(rw, resBytes, place)
	}

	printResult(err, http.StatusCreated, place)
	return
}
