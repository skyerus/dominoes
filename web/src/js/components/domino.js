import React, {useEffect} from "react"
import { makeStyles } from '@material-ui/core/styles'

const useStyles = makeStyles(theme => ({
  container: {
    border: "1px solid #fff",
    borderRadius: "2px",
    display: "flex",
  },
  horizontalContainer: {
    width: "160px",
    height: "80px",
  },
  verticalContainer: {
    width: "80px",
    height: "160px",
    flexDirection: "column",
  },
  numberContainer: {
    fontSize: "30px",
    display: "flex",
    justifyContent: "center",
    alignItems: "center",
    flexBasis: "50%",
  },
  verticalNumberContainer: {
    borderBottom: "1px solid #fff",
  },
  horizontalNumberContainer: {
    borderRight: "1px solid #fff",
  }
}))

export default function Domino(props) {
  const classes = useStyles()
  return (
    <div className={`${classes.container} ${props.vertical ? classes.verticalContainer : classes.horizontalContainer}`}>
      <div className={`${classes.numberContainer} ${props.vertical ? classes.verticalNumberContainer : classes.horizontalNumberContainer}`}>
        <div>{props.left}</div>
      </div>
      <div className={classes.numberContainer}>
        <div>{props.right}</div>
      </div>
    </div>
  )
}