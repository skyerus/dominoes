package main

import (
	"math/rand"
	"time"

	"github.com/skyerus/dominoes/pkg/api"
	"github.com/skyerus/dominoes/pkg/game"
)

func main() {
	// env.SetEnv()
	rand.Seed(time.Now().UnixNano())
	sessions := game.NewSessions()
	main := &api.App{}
	main.Initialize(sessions)
	main.Run(":8080")
}
