package main

import (
	api "db_forum/internal/services"
	utils "db_forum/internal/utils"
	"os"

	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	const confPath = "conf.json"

	API, conf, err := api.GetHandler(confPath)
	if err != nil {
		utils.PrintDebug("Error in configuration file or database" + err.Error())
		return
	}

	r := mux.NewRouter()

	r.HandleFunc("/api/user/{nickname}/create", API.UserCreate).Methods("POST")
	r.HandleFunc("/api/user/{nickname}/profile", API.UserUpdateProfile).Methods("POST")
	r.HandleFunc("/api/user/{nickname}/profile", API.UserProfile).Methods("GET")

	r.HandleFunc("/api/forum/create", API.ForumCreate).Methods("POST")
	r.HandleFunc("/api/forum/{slug}/create", API.ThreadCreate).Methods("POST")
	r.HandleFunc("/api/forum/{slug}/details", API.ForumDetails).Methods("GET")
	r.HandleFunc("/api/forum/{slug}/threads", API.ForumThreads).Methods("GET")
	r.HandleFunc("/api/forum/{slug}/users", API.ForumUsers).Methods("GET")

	r.HandleFunc("/api/thread/{slug_or_id}/create", API.ThreadCreatePost).Methods("POST")
	r.HandleFunc("/api/thread/{slug_or_id}/vote", API.ThreadVote).Methods("POST")
	r.HandleFunc("/api/thread/{slug_or_id}/details", API.ThreadUpdate).Methods("POST")
	r.HandleFunc("/api/thread/{slug_or_id}/details", API.ThreadDetails).Methods("GET")
	r.HandleFunc("/api/thread/{slug_or_id}/posts", API.ThreadPosts).Methods("GET")

	r.HandleFunc("/api/post/{id}/details", API.PostUpdate).Methods("POST")
	r.HandleFunc("/api/post/{id}/details", API.PostDetails).Methods("GET")

	r.HandleFunc("/api/service/status", API.ServiceStatus).Methods("GET")
	r.HandleFunc("/api/service/clear", API.ServiceClear).Methods("POST")

	utils.PrintDebug("launched, look at us on " + conf.Server.Host + ":" + conf.Server.Port)

	if os.Getenv("PORT") == "" {
		os.Setenv("PORT", conf.Server.Port)
	}

	if err = http.ListenAndServe(":"+os.Getenv("PORT"), r); err != nil {
		utils.PrintDebug("Error:" + err.Error())
	}
}
