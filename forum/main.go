package main

import (
	"time"
	"net/http/httputil"
	"fmt"
	api "db_forum/internal/services"
	utils "db_forum/internal/utils"
	"os"

	"net/http"

	"github.com/gorilla/mux"
)

var maxElapsed time.Duration = -1
var maxString = ""

func logHandler(fn http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
		// fn(w, r)
		// return
        x, err := httputil.DumpRequest(r, true)
        if err != nil {
            http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
            return
		}
		start := time.Now()
		fn(w, r)
		end := time.Now()
		elapsed := end.Sub(start)
		if elapsed > maxElapsed {
			maxElapsed = elapsed
			maxString = fmt.Sprintf("%q  --- duration = %s", x, elapsed.String())
			fmt.Println(maxString)
		}
		//fmt.Println(fmt.Sprintf("%q  --- duration = %s", x, elapsed.String())) 
    }
}

func main() {
	const confPath = "conf.json"

	API, conf, err := api.GetHandler(confPath)
	if err != nil {
		utils.PrintDebug("Error in configuration file or database" + err.Error())
		return
	}

	r := mux.NewRouter()

	r.HandleFunc("/api/user/{nickname}/create", logHandler(API.UserCreate)).Methods("POST")
	r.HandleFunc("/api/user/{nickname}/profile", logHandler(API.UserUpdateProfile)).Methods("POST")
	r.HandleFunc("/api/user/{nickname}/profile", logHandler(API.UserProfile)).Methods("GET")

	r.HandleFunc("/api/forum/create", logHandler(API.ForumCreate)).Methods("POST")
	r.HandleFunc("/api/forum/{slug}/create", logHandler(API.ThreadCreate)).Methods("POST")
	r.HandleFunc("/api/forum/{slug}/details", logHandler(API.ForumDetails)).Methods("GET")
	r.HandleFunc("/api/forum/{slug}/threads", logHandler(API.ForumThreads)).Methods("GET")
	r.HandleFunc("/api/forum/{slug}/users", logHandler(API.ForumUsers)).Methods("GET")

	r.HandleFunc("/api/thread/{slug_or_id}/create", logHandler(API.ThreadCreatePost)).Methods("POST")
	r.HandleFunc("/api/thread/{slug_or_id}/vote", logHandler(API.ThreadVote)).Methods("POST")
	r.HandleFunc("/api/thread/{slug_or_id}/details", logHandler(API.ThreadUpdate)).Methods("POST")
	r.HandleFunc("/api/thread/{slug_or_id}/details", logHandler(API.ThreadDetails)).Methods("GET")
	r.HandleFunc("/api/thread/{slug_or_id}/posts", logHandler(API.ThreadPosts)).Methods("GET")

	r.HandleFunc("/api/post/{id}/details", logHandler(API.PostUpdate)).Methods("POST")
	r.HandleFunc("/api/post/{id}/details", logHandler(API.PostDetails)).Methods("GET")

	r.HandleFunc("/api/service/status", logHandler(API.ServiceStatus)).Methods("GET")
	r.HandleFunc("/api/service/clear", logHandler(API.ServiceClear)).Methods("POST")

	utils.PrintDebug("launched, look at us on " + conf.Server.Host + ":" + conf.Server.Port)

	if os.Getenv("PORT") == "" {
		os.Setenv("PORT", conf.Server.Port)
	}

	if err = http.ListenAndServe(":"+os.Getenv("PORT"), r); err != nil {
		utils.PrintDebug("Error:" + err.Error())
	}
}
