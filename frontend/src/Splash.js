import React from 'react'
import { useState, useEffect, forwardRef } from "react";
import classNames from "classnames";
import { makeStyles } from "@material-ui/core/styles";

import Fade from "@material-ui/core/Fade";

import Box from '@material-ui/core/Box'
import Card from "@material-ui/core/Card";
import CardContent from "@material-ui/core/CardContent";
import CardHeader from "@material-ui/core/CardHeader";
import CardMedia from "@material-ui/core/CardMedia";
import Typography from "@material-ui/core/Typography";
import Grid from '@material-ui/core/Grid'

import useScrollTrigger from "@material-ui/core/useScrollTrigger";
import { AddBoxSharp } from '@material-ui/icons';

export const Splash = forwardRef((props, ref) => {

  const visited = useVisited();
  const wasScrolled = useScrolled();

  const timeout = visited ? 0 : 3000;
  const image = require('./big-pic-m.jpg')

  const useStyles = makeStyles((theme) => ({
    splash: {
      //postion: 'absolute',
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
      minHeight: "57vw",
      width: "100%",
      transition: "all .3s ease-in-out",
      transitionDelay: '.3s',
      //[theme.breakpoints.up("sm")]: {minHeight: '1vh'}
    },
    text: {
      [theme.breakpoints.down("xs")]: {fontSize: '50px'},
      //[theme.breakpoints.down("sm")]: {fontSize: '70px'},
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

  const classes = useStyles()

  return (
    <Card className={classNames({
          'classes.animateReveal': visited || (!visited && wasScrolled)
        })} 
        elevation={0}>
      <CardMedia ref={ref}>
        <img className={classes.media} src={image} alt='icon'/>
        <Grid
          container
          direction="column"
          className={classes.center}
        >
          <Box className={classes.centerVertical}>
            <Typography className={classes.text} color='primary' align='left' variant="h1">
              Full Stack.
              <Fade in={true} timeout={timeout}>
                <Box>
                  Big Picture.
                </Box>
              </Fade>
            </Typography>
            </Box>
          </Grid>
      </CardMedia>
    </Card>

      // {/* <CardMedia 
      //   height="400"
      //   image={image}
      // /> */}
      // {/* <img
      //   className={classNames("big-picture", {
      //     "animate-reveal": visited || (!visited && wasScrolled)
      //   })}
      //   ref={ref}
      //   src={require("./big_picture.jpg")}
      //   alt="icon"
      // /> */}
      // {/* <div className="center-content big-picture-text" style={{ height: 200 }}>

      // </div> */}
  );
});

const useScrolled = () => {
  const trigger = useScrollTrigger({ threshold: 0.1 });
  const [scrolled, setScrolled] = useState(false);

  useEffect(() => {
    if (trigger) {
      setScrolled(true);
    }
  }, [trigger]);

  return scrolled;
};

const useVisited = () => {
  const [value, setValue] = useState(false);

  useEffect(() => {
    setValue(getCookie("visited"));
    setCookie("visited", true);
  }, []);

  return value;
};

const getCookie = (cookie) => {
  let value = document.cookie.match("(^|;)\\s*" + cookie + "\\s*=\\s*([^;]+)");
  return value ? value.pop() : "";
};

const setCookie = (cookie, value) => {
  document.cookie = `${cookie}=${value}`;
};