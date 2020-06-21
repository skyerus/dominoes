import React, {useEffect} from "react"
import Domino from './domino'

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
    <div>
      <div style={{position: "absolute", bottom: "0", marginLeft: "20px", marginRight: "20px"}}>
        <Domino left={0} right={1} vertical={true}/>
      </div>
    </div>
  )
}
