import React, {useEffect} from "react"
import { makeStyles } from '@material-ui/core/styles'
import FormControl from '@material-ui/core/FormControl'
import InputLabel from '@material-ui/core/InputLabel'
import Select from '@material-ui/core/Select'
import MenuItem from '@material-ui/core/MenuItem'
import Button from '@material-ui/core/Button'
import Dialog from '@material-ui/core/Dialog'
import DialogActions from '@material-ui/core/DialogActions'
import DialogContent from '@material-ui/core/DialogContent'
import DialogTitle from '@material-ui/core/DialogTitle'

import Domino from './domino'

import coreApi from '../api/core'

const useStyles = makeStyles(theme => ({
  playerTilesContainer: {
    backgroundColor: "#666666",
    position: "absolute", 
    bottom: "0", 
    width: "100%",
    display: "flex",
    justifyContent: "center",
    paddingTop: "20px",
    paddingBottom: "20px",
  },
  playerDominoContainer: {
    marginLeft: "7px",
    marginRight: "7px",
  },
  boardContainer: {
    marginTop: "20%",
    display: "flex",
    justifyItems: "center",
    overflowX: "auto",
    marginLeft: "20px",
    marginRight: "20px",
    paddingBottom: "4px",
  },
  boardDominoContainer: {
    marginLeft: "3px",
    marginRight: "3px",
  },
  drawContainer: {
    marginTop: "20px",
    display: "flex",
    justifyContent: "center",
  },
  draw: {
    width: "120px",
    height: "120px",
    borderRadius: "50%",
    border: "2px solid #fff",
    fontSize: "30px",
    color: "#fff",
    textAlign: "center",
    lineHeight: "120px",
  },
  counter: {
    fontSize: "20px",
    textAlign: "center",
    display: "flex",
    alignItems: "center",
    paddingLeft: "10px",
    paddingRight: "10px",
  },
  newGame: {
    display: "flex",
    flexDirection: "column",
    justifyContent: "center",
    alignItems: "center",
  },
  formControl: {
    margin: theme.spacing(1),
    minWidth: 300,
  },
}))

export default function Game(props) {
  const classes = useStyles()
  const [gameState, setGameState] = React.useState(null)
  const [numOfPlayers, setNumOfPlayers] = React.useState(2)
  const [gameOver, setGameOver] = React.useState(false)

  useEffect(() => {
    fetchSession()
  }, [])

  const fetchSession = () => {
    coreApi.getSession().then(res => {
      updateGameState(res.data)
    })
  }

  const newSession = () => {
    coreApi.newSession(numOfPlayers).then(res => {
      updateGameState(res.data)
    })
  }

  const playTurn = (index) => {
    coreApi.playTurn(index).then(res => {
      updateGameState(res.data)
    }).catch(e => {
      props.snack("error", e.response.data.message)
    })
  }

  const updateGameState = (session) => {
    setGameState(session)
    setGameOver(session.gameover)
  }

  const draw = () => {
    coreApi.draw().then(res => {
      updateGameState(res.data)
    }).catch(e => {
      props.snack("error", e.response.data.message)
    })
  }

  const handleSelectChange = (e) => {
    setNumOfPlayers(e.target.value)
  }

  const closeWinModal = () => {
    setGameOver(false)
    setGameState(null)
  }

  return gameState ? (
    <div>
      <Dialog open={gameOver} onClose={() => closeWinModal} aria-labelledby="form-dialog-title">
        <DialogTitle id="form-dialog-title">{gameState.player_wins ? "You win!" : "Computer wins"}</DialogTitle>
        <DialogContent>
          {gameState.player_wins ? "Congratulations" : "Better luck next time"}
        </DialogContent>
        <DialogActions>
          <Button onClick={closeWinModal} color="primary">
            Close
          </Button>
      </DialogActions>
      </Dialog>
      <div className={classes.boardContainer}>
        {
          gameState.played_tiles.map((t, i) =>
            <div className={classes.boardDominoContainer}>
              <Domino key={i}
                      left={t.left} 
                      right={t.right}
                      vertical={false}
              />
            </div>
          )
        }
      </div>
      <div className={classes.drawContainer}>
        <div className={classes.counter}>
          Remaining tiles: {gameState.remaining_tiles}
        </div>
        <div className={classes.draw} onClick={draw}>
          DRAW
        </div>
      </div>
      <div className={classes.playerTilesContainer}>
        {
          gameState.my_tiles.map((t, i) =>
            <div className={classes.playerDominoContainer} onClick={() => playTurn(i)}>
              <Domino key={i}
                      left={t.left} 
                      right={t.right} 
                      vertical={true}
              />
            </div>
          )
        }
      </div>
    </div>
  ) : (
    <div className={classes.newGame}>
      <FormControl className={classes.formControl}>
        <InputLabel id="elect-label">Number of players</InputLabel>
        <Select
          labelId="select-label"
          id="simple-select"
          value={numOfPlayers}
          onChange={handleSelectChange}
        >
          <MenuItem value={2}>2</MenuItem>
          <MenuItem value={3}>3</MenuItem>
          <MenuItem value={4}>4</MenuItem>
        </Select>
      </FormControl>
      <Button onClick={newSession} color="primary">
        New game
      </Button>
    </div>
  )
}
