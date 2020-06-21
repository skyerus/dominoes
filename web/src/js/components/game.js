import React, {useEffect} from "react"

import coreApi from '../api/core'

export default function Game(props) {
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

  return (
    <div>GAME</div>
  )
}
