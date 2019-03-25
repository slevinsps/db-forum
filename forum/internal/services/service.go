package api

import (
	"db_forum/internal/models"
	"net/http"
)

//ServiceStatus
func (h *Handler) ServiceStatus(rw http.ResponseWriter, r *http.Request) {

	const place = "ServiceStatus"
	var (
		err         error
		countPost   int
		countUser   int
		countThread int
		countForum  int
	)

	if countPost, err = h.DB.CountPost(); err != nil {
		rw.WriteHeader(http.StatusNotFound)
		printResult(err, http.StatusNotFound, place)
		return
	}
	if countForum, err = h.DB.CountForum(); err != nil {
		rw.WriteHeader(http.StatusNotFound)
		printResult(err, http.StatusNotFound, place)
		return
	}
	if countThread, err = h.DB.CountThread(); err != nil {
		rw.WriteHeader(http.StatusNotFound)
		printResult(err, http.StatusNotFound, place)
		return
	}
	if countUser, err = h.DB.CountUser(); err != nil {
		rw.WriteHeader(http.StatusNotFound)
		printResult(err, http.StatusNotFound, place)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	sendJSON(rw, models.Service{Forum: countForum, Post: countPost, User: countUser, Thread: countThread}, place)
	printResult(err, http.StatusCreated, place)
	return
}

//ServiceClear
func (h *Handler) ServiceClear(rw http.ResponseWriter, r *http.Request) {

	const place = "ServiceStatus"
	var (
		err error
	)

	if err = h.DB.ServiceClear(); err != nil {
		rw.WriteHeader(http.StatusNotFound)
		printResult(err, http.StatusNotFound, place)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	printResult(err, http.StatusCreated, place)
	return
}
