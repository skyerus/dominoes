package api

import (
	"net/http"
	"os"
	"strconv"
	"time"
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
