import React from 'react'
import { forwardRef } from "react";

import { makeStyles } from "@material-ui/core/styles";

import Box from '@material-ui/core/Box'
import Card from "@material-ui/core/Card";
import CardMedia from "@material-ui/core/CardMedia";
import Fade from "@material-ui/core/Fade";
import Grid from '@material-ui/core/Grid'
import Typography from "@material-ui/core/Typography";

import DynamicPicture from './DynamicPicture'

const Splash = forwardRef((props, ref) => {
  const classes = useStyles()

  return (
    <Card elevation={0}>
      <CardMedia ref={ref}>
        <Box className={classes.mediaWrapper} >
          <DynamicPicture className={classes.media} src="https://cmelgreen.com/static/media/big-pic.jpg" />
        </Box>
        <Grid container direction="column" className={classes.center}>
          <Box className={classes.centerVertical}>
            <Typography className={classes.text} color='primary' align='left' variant="h1">
              Full Stack.
              <Fade in={true} timeout={props.timeout}>
                <Box>
                  Big Picture.
                </Box>
              </Fade>
            </Typography>
            </Box>
          </Grid>
      </CardMedia>
    </Card>


  );
});

const useStyles = makeStyles((theme) => ({
  splash: {
    display: 'grid',
    gridTemplateRows: '1fr',
    gridTemplateColumns: '1fr',
    margin: 0,
    padding: 0,
    border: 0,
  },
  media: {
    top: 0,
    right: 0,
    minHeight: "100%",
    width: "100%",
    transition: "all .3s ease-in-out",
    transitionDelay: '.3s',
  },
  mediaWrapper: {
    height: '57vw',
    width: '100vw'
  },
  text: {
    [theme.breakpoints.down("xs")]: {fontSize: '50px'},
  }, 
  center: {
    position: 'absolute',
    height: '100%',
    top: 0,
    right: 0,
    alignItems: 'center',
    alignContent: 'center',
    justify: 'center',
  },
  centerVertical: {
    marginTop: '16%',
    position: 'relative',
    height: 'inheret'
  }
}))

export default Splash