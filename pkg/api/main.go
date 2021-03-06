package api

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/skyerus/dominoes/pkg/game"
)

// App - Root component
type App struct {
	Router *mux.Router
}

type router struct {
	sessions *game.Sessions
}

func newRouter(sessions *game.Sessions) *router {
	return &router{sessions}
}

// Initialize - Initialize app
func (a *App) Initialize(sessions *game.Sessions) {
	router := newRouter(sessions)
	a.Router = mux.NewRouter()
	// a.Router.Use(cors)
	a.setRouters(router)
}

func (a *App) setRouters(router *router) {
	// Base routes
	a.Router.HandleFunc("/api/session", router.currentSession).Methods("GET", "OPTIONS")
	a.Router.HandleFunc("/api/new_game", router.newGame).Methods("POST", "OPTIONS")
	a.Router.HandleFunc("/api/play_turn/{index}", router.playTurn).Methods("POST", "OPTIONS")
	a.Router.HandleFunc("/api/draw", router.drawTile).Methods("POST", "OPTIONS")
	a.Router.PathPrefix("/").Handler(http.FileServer(http.Dir("./web/dist")))
}

// Run - Run the app
func (a *App) Run(host string) {
	srv := &http.Server{
		Handler:      a.Router,
		Addr:         host,
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  18 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
}
