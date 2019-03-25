package api

import (
	"db_forum/internal/models"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
)

func (h *Handler) getNickname(r *http.Request) (nickname string, err error) {
	var (
		vars map[string]string
	)

	vars = mux.Vars(r)

	if nickname = vars["nickname"]; nickname == "" {
		err = errors.New("Cant found parameters")
		return
	}
	return
}

func (h *Handler) getId(r *http.Request) (id string, err error) {
	var (
		vars map[string]string
	)

	vars = mux.Vars(r)

	if id = vars["id"]; id == "" {
		err = errors.New("Cant found parameters")
		return
	}
	return
}

func (h *Handler) getSlug(r *http.Request) (slug string, err error) {
	var (
		vars map[string]string
	)

	vars = mux.Vars(r)

	if slug = vars["slug"]; slug == "" {
		err = errors.New("Cant found parameters")
		return
	}
	return
}

func (h *Handler) getSlugOrID(r *http.Request) (slugOrId string, err error) {
	var (
		vars map[string]string
	)

	vars = mux.Vars(r)

	if slugOrId = vars["slug_or_id"]; slugOrId == "" {
		err = errors.New("Cant found parameters")
		return
	}
	return
}

func getUserFromBody(r *http.Request) (user models.User, emptyBody bool, err error) {

	emptyBody = false
	if r.Body == nil {
		err = errors.New("Cant found parameters")

		return
	}
	var body []byte
	body, err = ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if len(body) == 3 {
		emptyBody = true
	}
	_ = json.Unmarshal(body, &user)

	return
}

func getVoteFromBody(r *http.Request) (vote models.Vote, err error) {

	if r.Body == nil {
		err = errors.New("Cant found parameters")

		return
	}
	var body []byte
	body, err = ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	_ = json.Unmarshal(body, &vote)

	return
}

func getForumFromBody(r *http.Request) (forum models.Forum, err error) {

	if r.Body == nil {
		err = errors.New("Cant found parameters")

		return
	}
	defer r.Body.Close()

	_ = json.NewDecoder(r.Body).Decode(&forum)

	return
}

func getThreadFromBody(r *http.Request) (thread models.Thread, err error) {

	if r.Body == nil {
		err = errors.New("Cant found parameters")

		return
	}
	defer r.Body.Close()

	_ = json.NewDecoder(r.Body).Decode(&thread)

	return
}

func getMessageFromBody(r *http.Request) (message models.Message, err error) {

	if r.Body == nil {
		err = errors.New("Cant found parameters")

		return
	}
	defer r.Body.Close()

	_ = json.NewDecoder(r.Body).Decode(&message)

	return
}

func getPostsFromBody(r *http.Request) (post []models.Post, err error) {

	if r.Body == nil {
		err = errors.New("Cant found parameters")

		return
	}
	defer r.Body.Close()

	_ = json.NewDecoder(r.Body).Decode(&post)

	return
}
