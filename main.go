package main

import (
	"github.com/skyerus/dominoes/pkg/api"
	"github.com/skyerus/dominoes/pkg/game"
)

func main() {
	// env.SetEnv()
	sessions := game.NewSessions()
	main := &api.App{}
	main.Initialize(sessions)
	main.Run(":8080")
}
