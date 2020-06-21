import React, {useEffect} from "react"

import { createMuiTheme } from '@material-ui/core/styles'
import { ThemeProvider } from '@material-ui/styles'
import Snackbar from '@material-ui/core/Snackbar'
import MuiAlert from '@material-ui/lab/Alert'
import CssBaseline from '@material-ui/core/CssBaseline'

import Game from './components/game'

function getTheme(theme) {
  return createMuiTheme({
    palette: {
      type: theme.paletteType,
      primary: {
        main: theme.paletteType === 'light' ? '#1976d2' : '#3f51b5',
      },
      secondary: {
        main: '#03DAC5',
      },
      contrastThreshold: 3,
      tonalOffset: 0.2,
      background: {
        default: theme.paletteType === 'light' ? '#f7f9fc' : '#303030'
      },
    },
  })
}

function Alert(props) {
  return <MuiAlert elevation={6} variant="filled" {...props} />
}

export default function App(props) {
  const [state, setState] = React.useState({
    success: {
      active: false,
      message: ""
    },
    error: {
      active: false,
      message: ""
    },
  })

  let theme = getTheme({
    paletteType: 'dark',
  })
  
  const handleClose = (e, reason, type) => {
    if (reason === 'clickaway') {
      return
    }
    setState({...state, [type]: {active: false, message: ""}})
  }

  const snack = (type, message) => {
    setState({...state, [type]: {message: message, active: true}})
  }

  const commonProps = {
    snack: snack,
  }


  return (
    <ThemeProvider theme={theme}>
      <Snackbar open={state.success.active} autoHideDuration={3000} onClose={(e, reason) => handleClose(e, reason, "success")}>
        <Alert severity="success" onClose={(e, reason) => handleClose(e, reason, "success")}>
          {state.success.message}
        </Alert>
      </Snackbar>
      <Snackbar open={state.error.active} autoHideDuration={3000} onClose={(e, reason) => handleClose(e, reason, "error")}>
        <Alert severity="error" onClose={(e, reason) => handleClose(e, reason, "error")}>
          {state.error.message}
        </Alert>
      </Snackbar>
      <Game {...commonProps}/>
      <CssBaseline/>
    </ThemeProvider>
  )
}
