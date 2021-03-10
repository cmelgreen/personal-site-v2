import React from 'react'
import { forwardRef } from 'react'

import { makeStyles } from "@material-ui/core/styles";
import { Container, Typography } from '@material-ui/core'

const useStyles = makeStyles((theme) => ({
  aboutMe: {
    backgroundColor: '#2C3359',
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

const About = forwardRef((props, ref) => {
  const classes = useStyles()

  return (
    <Container className={classes.aboutMe} maxWidth='md' ref={ref}>
        <Typography className={classes.title} variant="h1" align="center">
            About
        </Typography>
        <Typography className={classes.text} variant="subtitle1">
            After 6 great years of Enterprise Software-as-a-Service sales, pivoting back towards the technical side of SaaS and full-stack software development
        </Typography>
    </Container>
  )
})

export default About