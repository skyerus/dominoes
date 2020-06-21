package api

import (
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

func (router *router) newGame(w http.ResponseWriter, r *http.Request) {
	numOfPlayers, err := strconv.Atoi(r.URL.Query().Get("numOfPlayers"))
	if err != nil {
		respondBadRequest(w)
		return
	}
	sessionID := strconv.Itoa(int(time.Now().UnixNano()))
	session, customErr := router.sessions.NewSession(sessionID, numOfPlayers)
	if customErr != nil {
		handleError(w, customErr)
		return
	}
	cookie := &http.Cookie{Name: "session-id", Value: sessionID, Domain: os.Getenv("API_DOMAIN"), MaxAge: 7200, Path: "/"}
	http.SetCookie(w, cookie)
	respondJSON(w, http.StatusOK, session)
}

func (router *router) playTurn(w http.ResponseWriter, r *http.Request) {
	sessionCookie, err := r.Cookie("session-id")
	if err != nil {
		respondBadRequest(w)
		return
	}
	indexStr, success := mux.Vars(r)["index"]
	if !success {
		respondBadRequest(w)
		return
	}
	index, err := strconv.Atoi(indexStr)
	if err != nil {
		respondBadRequest(w)
		return
	}
	sessionID := sessionCookie.Value
	session := router.sessions.FetchSession(sessionID)
	if session == nil {
		respondBadRequest(w)
		return
	}
	customErr := session.PlayTurn(index)
	if customErr != nil {
		handleError(w, customErr)
		return
	}

	respondJSON(w, http.StatusOK, session)
}
