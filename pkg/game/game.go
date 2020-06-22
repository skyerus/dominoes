package game

import (
	"errors"
	"math"
	"math/rand"
	"strconv"
	"sync"

	"github.com/skyerus/dominoes/pkg/customerror"
)

// Session - game session state
type Session struct {
	mux            sync.Mutex
	Players        []*player `json:"players"`
	PlayedTiles    *[]tile   `json:"played_tiles"`
	RemainingTiles *[]tile   `json:"remaining_tiles"`
	PlayersTurn    int       `json:"players_turn"`
	Gameover       bool      `json:"gameover"`
	Playerwins     bool      `json:"player_wins"`
}

// FormattedSession - for front-end
type FormattedSession struct {
	MyTiles        []tile `json:"my_tiles"`
	Gameover       bool   `json:"gameover"`
	Playerwins     bool   `json:"player_wins"`
	PlayedTiles    []tile `json:"played_tiles"`
	RemainingTiles int    `json:"remaining_tiles"`
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
	rand.Shuffle(numOfTiles(maxNumOfPips), func(i, j int) { (*tiles)[i], (*tiles)[j] = (*tiles)[j], (*tiles)[i] })
	players := make([]*player, numOfPlayers)
	var realPlayer player
	playerTiles := (*tiles)[0:tilesPerPlayer]
	realPlayer.Tiles = &playerTiles
	players[0] = &realPlayer
	for i := 1; i < numOfPlayers; i++ {
		var p player
		p.IsBot = true
		playerTiles := (*tiles)[i*tilesPerPlayer : ((i + 1) * tilesPerPlayer)]
		p.Tiles = &playerTiles
		players[i] = &p
	}
	remainingTiles := (*tiles)[0:(numOfTiles(maxNumOfPips) - tilesPerPlayer*numOfPlayers)]
	doubleFound := orderPlayers(&players)
	if !doubleFound {
		return NewSession(numOfPlayers)
	}
	playedTiles := make([]tile, 0)
	s.PlayedTiles = &playedTiles
	s.Players = players
	s.RemainingTiles = &remainingTiles

	s.preGame()
	return &s, nil
}

// PlayTurn - user plays their turn
func (s *Session) PlayTurn(tileIndex int) customerror.Error {
	s.mux.Lock()
	defer s.mux.Unlock()
	if s.Players[s.PlayersTurn].IsBot {
		return customerror.NewBadRequestError("Not your turn")
	}
	if tileIndex > len(*s.Players[s.PlayersTurn].Tiles)-1 {
		return customerror.NewBadRequestError("Invalid tile")
	}
	t := (*s.Players[s.PlayersTurn].Tiles)[tileIndex]
	err := s.placeTile(t)
	if err != nil {
		return customerror.NewBadRequestError("Invalid tile")
	}
	*s.Players[s.PlayersTurn].Tiles = removeTile(*s.Players[s.PlayersTurn].Tiles, tileIndex)
	s.incrementTurn()
	s.playBotTurns()
	return nil
}

// DrawTile - player draws tile
func (s *Session) DrawTile() customerror.Error {
	s.mux.Lock()
	defer s.mux.Unlock()
	if s.Players[s.PlayersTurn].IsBot {
		return customerror.NewBadRequestError("Not your turn")
	}
	err := s.drawTile(s.Players[s.PlayersTurn])
	if err != nil {
		s.endGame()
		return nil
	}
	s.incrementTurn()
	s.playBotTurns()
	return nil
}

// FormatSession - format session for front
func FormatSession(s *Session) FormattedSession {
	var fSession FormattedSession
	var me *player
	for _, p := range s.Players {
		if !p.IsBot {
			me = p
			break
		}
	}
	fSession.MyTiles = *me.Tiles
	fSession.Gameover = s.Gameover
	fSession.Playerwins = s.Playerwins
	fSession.PlayedTiles = *s.PlayedTiles
	fSession.RemainingTiles = len(*s.RemainingTiles)

	return fSession
}

func (s *Session) endGame() {
	s.Gameover = true
	for _, p := range s.Players {
		if len(*p.Tiles) == 0 {
			if !p.IsBot {
				s.Playerwins = true
				return
			}
			return
		}
	}
	var winner player
	var lowestScore int
	for _, p := range s.Players {
		var score int
		for _, t := range *p.Tiles {
			score += t.Left + t.Right
		}
		if score < lowestScore || lowestScore == 0 {
			lowestScore = score
			winner = *p
		}
	}
	if !winner.IsBot {
		s.Playerwins = true
	}
}

func (s *Session) botTurn(p *player) error {
	for i, t := range *p.Tiles {
		err := s.placeTile(t)
		if err == nil {
			*p.Tiles = removeTile(*p.Tiles, i)
			s.incrementTurn()
			return nil
		}
	}
	err := s.drawTile(p)
	s.incrementTurn()

	return err
}

func (s *Session) incrementTurn() {
	if s.PlayersTurn == len(s.Players)-1 {
		s.PlayersTurn = 0
		return
	}
	s.PlayersTurn++
}

func (s *Session) placeTile(t tile) error {
	leftTile := (*s.PlayedTiles)[0]
	rightTile := (*s.PlayedTiles)[len(*s.PlayedTiles)-1]
	var placeLeft bool
	if t.Left == leftTile.Left {
		placeLeft = true
		t.Left, t.Right = t.Right, t.Left
	} else if t.Right == leftTile.Left {
		placeLeft = true
	} else if t.Left == rightTile.Right {
		placeLeft = false
	} else if t.Right == rightTile.Right {
		placeLeft = false
		t.Left, t.Right = t.Right, t.Left
	} else {
		return errors.New("Illegal move")
	}

	if placeLeft {
		*s.PlayedTiles = append([]tile{t}, *s.PlayedTiles...)
	} else {
		*s.PlayedTiles = append(*s.PlayedTiles, t)
	}

	return nil
}

func (s *Session) drawTile(p *player) error {
	if len(*s.RemainingTiles) == 0 {
		return errors.New("Out of tiles")
	}
	randIndex := rand.Intn(len(*s.RemainingTiles))
	randTile := (*s.RemainingTiles)[randIndex]
	*p.Tiles = append(*p.Tiles, randTile)
	*s.RemainingTiles = removeTile(*s.RemainingTiles, randIndex)

	return nil
}

func (s *Session) playBotTurns() {
	for {
		if !s.Players[s.PlayersTurn].IsBot {
			return
		}
		err := s.botTurn(s.Players[s.PlayersTurn])
		if err != nil {
			s.endGame()
			return
		}
	}
}

func (s *Session) preGame() {
	s.playHighestDouble((s.Players)[0])
	s.incrementTurn()
	for s.PlayersTurn != 0 {
		if !(*s.Players[s.PlayersTurn]).IsBot {
			break
		}
		err := s.botTurn(s.Players[s.PlayersTurn])
		if err != nil {
			s.endGame()
			return
		}
	}
}

func (s *Session) playHighestDouble(p *player) {
	highestDouble := tile{Left: -1, Right: -1}
	var iOfHighestDouble int
	for i, t := range *p.Tiles {
		if t.Left == t.Right {
			if t.Left > highestDouble.Left {
				iOfHighestDouble = i
				highestDouble = t
			}
		}
	}
	*p.Tiles = removeTile(*p.Tiles, iOfHighestDouble)
	*s.PlayedTiles = append(*s.PlayedTiles, highestDouble)
}

func removeTile(ts []tile, index int) []tile {
	return append(ts[:index], ts[index+1:]...)
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

func orderPlayers(players *[]*player) bool {
	iPlayerOne, doubleFound := indexOfPlayerOne(players)
	if !doubleFound {
		return doubleFound
	}
	(*players)[0], (*players)[iPlayerOne] = (*players)[iPlayerOne], (*players)[0]
	return true
}

func indexOfPlayerOne(players *[]*player) (int, bool) {
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
