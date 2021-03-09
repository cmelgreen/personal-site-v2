import React from 'react'

import { makeStyles } from "@material-ui/core/styles";
import { Container, Typography } from '@material-ui/core'

const useStyles = makeStyles((theme) => ({
  aboutMe: {
    backgroundColor: '#353440',
    padding: '10%',
    margin: '5% auto',
    minHeight: 200,
    [theme.breakpoints.down("sm")]: {margin: '0 auto'},
  },
  title: {
    color: 'white',
    [theme.breakpoints.down("xs")]: {fontSize: '50px'},
  },
  text: {
    color: 'white',
  },
}))

export default function About(props) {
  const classes = useStyles()

  return (
    <Container className={classes.aboutMe} maxWidth='md'>
      <Typography className={classes.title} variant="h1" align="center">
        About Me
      </Typography>
      <Typography className={classes.text} variant="subtitle1">
      "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum."
      </Typography>
    </Container>
  )
}