package api

import (
	"encoding/json"
	"db_forum/internal/models"
	"net/http"
)

func (h *Handler) UserCreate(rw http.ResponseWriter, r *http.Request) {
	const place = "UserCreate"

	var (
		err         error
		users       []models.User
		nickname    string
		user        models.User
		checkUnique bool
	)

	if nickname, err = h.getNickname(r); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		printResult(err, http.StatusBadRequest, place)
		return
	}

	if user, _, err = getUserFromBody(r); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		printResult(err, http.StatusBadRequest, place)
		return
	}

	user.Nickname = nickname

	if users, checkUnique, err = h.DB.CreateUser(user); err != nil {
		rw.WriteHeader(http.StatusNotFound)
		printResult(err, http.StatusNotFound, place)
		return
	}
	rw.Header().Set("Content-Type", "application/json")

	if checkUnique {
		rw.WriteHeader(http.StatusCreated)
		resBytes, _ := users[0].MarshalJSON()
		sendJSON(rw, resBytes, place)
		
	} else {
		rw.WriteHeader(http.StatusConflict)

		resBytes, _ := json.Marshal(users)
		sendJSON(rw, resBytes, place)
		
	}

	printResult(err, http.StatusCreated, place)
	return
}

func (h *Handler) UserProfile(rw http.ResponseWriter, r *http.Request) {
	const place = "UserProfile"

	var (
		err           error
		nickname      string
		user          models.User
		checkFindUser bool
	)

	if nickname, err = h.getNickname(r); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		printResult(err, http.StatusBadRequest, place)
		return
	}

	if user, checkFindUser, err = h.DB.GetUserByNickname(nickname); err != nil {
		rw.WriteHeader(http.StatusNotFound)
		printResult(err, http.StatusNotFound, place)
		return
	}

	rw.Header().Set("Content-Type", "application/json")

	if checkFindUser {
		rw.WriteHeader(http.StatusOK)
		resBytes, _ := user.MarshalJSON()
		sendJSON(rw, resBytes, place)
	} else {
		rw.WriteHeader(http.StatusNotFound)
		message := models.Message{Message: "Can't find user by nickname: " + nickname}

		resBytes, _ := message.MarshalJSON()
		sendJSON(rw, resBytes, place)
	}

	printResult(err, http.StatusCreated, place)
	return
}

func (h *Handler) UserUpdateProfile(rw http.ResponseWriter, r *http.Request) {
	const place = "UserUpdateProfile"

	var (
		err           error
		nickname      string
		user          models.User
		userQuery     models.User
		checkUnique   bool
		checkFindUser bool
		//checkEmptyBody bool
	)

	if nickname, err = h.getNickname(r); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		printResult(err, http.StatusBadRequest, place)
		return
	}

	if user, _, err = getUserFromBody(r); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		printResult(err, http.StatusBadRequest, place)
		return
	}

	rw.Header().Set("Content-Type", "application/json")

	user.Nickname = nickname

	if userQuery, checkFindUser, err = h.DB.GetUserByNickname(nickname); err != nil {
		rw.WriteHeader(http.StatusNotFound)
		printResult(err, http.StatusNotFound, place)
		return
	}

	if !checkFindUser {
		rw.WriteHeader(http.StatusNotFound)
		message := models.Message{Message: "Can't find user by nickname: " + nickname}

		resBytes, _ := message.MarshalJSON()
		sendJSON(rw, resBytes, place)
		return
	}

	if user.About == "" {
		user.About = userQuery.About
	}
	if user.Email == "" {
		user.Email = userQuery.Email
	}
	if user.Fullname == "" {
		user.Fullname = userQuery.Fullname
	}
	if user.Nickname == "" {
		user.Nickname = userQuery.Nickname
	}

	if checkUnique, err = h.DB.UpdateUser(user); err != nil {

		rw.WriteHeader(http.StatusNotFound)
		printResult(err, http.StatusNotFound, place)
		return
	}

	if checkUnique {
		rw.WriteHeader(http.StatusOK)
		resBytes, _ := user.MarshalJSON()
		sendJSON(rw, resBytes, place)
	} else {
		rw.WriteHeader(http.StatusConflict)
		message := models.Message{Message: "Can't find user by nickname: " + nickname}

		resBytes, _ := message.MarshalJSON()
		sendJSON(rw, resBytes, place)
	}

	printResult(err, http.StatusCreated, place)
	return
}

func (h *Handler) ForumUsers(rw http.ResponseWriter, r *http.Request) {

	const place = "ForumUsers"

	var (
		err            error
		slug           string
		forum          models.Forum
		users          []models.User
		checkFindForum bool
	)

	if slug, err = h.getSlug(r); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		printResult(err, http.StatusBadRequest, place)
		return
	}
	rw.Header().Set("Content-Type", "application/json")

	if forum, checkFindForum, err = h.DB.GetForumBySlug(slug); err != nil {
		rw.WriteHeader(http.StatusNotFound)
		printResult(err, http.StatusNotFound, place)
		return
	}

	if !checkFindForum {
		rw.WriteHeader(http.StatusNotFound)
		message := models.Message{Message: "Can't find forum by sluq: " + slug}

		resBytes, _ := message.MarshalJSON()
		sendJSON(rw, resBytes, place)
		return
	}

	if users, err = h.DB.GetUsersByForum(forum.Slug, r.FormValue("limit"), r.FormValue("since"), r.FormValue("desc")); err != nil {
		rw.WriteHeader(http.StatusNotFound)
		printResult(err, http.StatusNotFound, place)
		return
	}
	rw.WriteHeader(http.StatusOK)
	if len(users) == 0 {
		rw.Write([]byte("[]"))
	} else {

		resBytes, _ := json.Marshal(users)
		sendJSON(rw, resBytes, place)
	}

	printResult(err, http.StatusCreated, place)
	return
}
