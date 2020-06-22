import React, {useEffect} from "react"
import { makeStyles } from '@material-ui/core/styles'
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
  }
}))

export default function Game(props) {
  const classes = useStyles()
  const [gameState, setGameState] = React.useState(null)

  useEffect(() => {
    fetchSession()
  }, [])

  const fetchSession = () => {
    coreApi.getSession().then(res => {
      setGameState(res.data)
    }).catch(e => {
      if (e.response.status == 400) {
        newSession()
      } else {
        props.snack("error", e.response.data.message)
      }
    })
  }

  const newSession = () => {
    coreApi.newSession(4).then(res => {
      setGameState(res.data)
    })
  }

  const playTurn = (index) => {
    coreApi.playTurn(index).then(res => {
      setGameState(res.data)
    }).catch(e => {
      props.snack("error", e.response.data.message)
    })
  }

  return gameState && (
    <div>
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
  )
}
