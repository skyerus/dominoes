package game

import (
	"math"
	"math/rand"
	"strconv"
	"time"

	"github.com/skyerus/dominoes/pkg/customerror"
)

// Session - game session state
type Session struct {
	Players        *[]player `json:"players"`
	PlayedTiles    *[]tile   `json:"played_tiles"`
	RemainingTiles *[]tile   `json:"remaining_tiles"`
}

type player struct {
	IsBot bool    `json:"isBot"`
	Tiles *[]tile `json:"tiles"`
}

type tile struct {
	Left  int `json:"left"`
	Right int `json:"right"`
}

const maxNumOfPips = 6
const maxPlayers = 4
const tilesPerPlayer = 7

// NewSession - new game session
func NewSession(numOfPlayers int) (*Session, customerror.Error) {
	var s Session
	if numOfPlayers > maxPlayers {
		return nil, customerror.NewBadRequestError("Max number of players is " + strconv.Itoa(maxPlayers))
	}
	tiles := allTiles(maxNumOfPips)
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(numOfTiles(maxNumOfPips), func(i, j int) { (*tiles)[i], (*tiles)[j] = (*tiles)[j], (*tiles)[i] })
	players := make([]player, numOfPlayers)
	var realPlayer player
	playerTiles := (*tiles)[0:tilesPerPlayer]
	realPlayer.Tiles = &playerTiles
	players[0] = realPlayer
	for i := 1; i < numOfPlayers; i++ {
		var p player
		p.IsBot = true
		playerTiles := (*tiles)[i*tilesPerPlayer : ((i + 1) * tilesPerPlayer)]
		p.Tiles = &playerTiles
		players[i] = p
	}
	remainingTiles := (*tiles)[0:(numOfTiles(maxNumOfPips) - tilesPerPlayer*numOfPlayers)]
	doubleFound := orderPlayers(&players)
	if !doubleFound {
		return NewSession(numOfPlayers)
	}
	playedTiles := make([]tile, 0)
	s.PlayedTiles = &playedTiles
	s.Players = &players
	s.RemainingTiles = &remainingTiles

	return &s, nil
}

func allTiles(maxNumOfPips int) *[]tile {
	tiles := make([]tile, numOfTiles(maxNumOfPips))
	index := 0
	for i := 0; i <= maxNumOfPips; i++ {
		for j := i; j <= maxNumOfPips; j++ {
			tiles[index] = tile{Left: i, Right: j}
			index++
		}
	}

	return &tiles
}

func numOfTiles(maxNumOfPips int) int {
	return (int(math.Pow(float64(maxNumOfPips), 2)) + 3*maxNumOfPips + 2) / 2
}

func orderPlayers(players *[]player) bool {
	iPlayerOne, doubleFound := indexOfPlayerOne(players)
	if !doubleFound {
		return doubleFound
	}
	(*players)[0], (*players)[iPlayerOne] = (*players)[iPlayerOne], (*players)[0]
	return true
}

func indexOfPlayerOne(players *[]player) (int, bool) {
	highestDouble := tile{Left: -1, Right: -1}
	var doubleFound bool
	iOfPlayerOne := 0
	for i, p := range *players {
		for _, t := range *p.Tiles {
			if t.Left == t.Right {
				if t.Left == maxNumOfPips {
					return i, true
				}
				if t.Left > highestDouble.Left {
					doubleFound = true
					iOfPlayerOne = i
					highestDouble = t
				}
			}
		}
	}

	return iOfPlayerOne, doubleFound
}
