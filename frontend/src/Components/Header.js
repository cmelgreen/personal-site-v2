import React from 'react'
import { useState, useEffect, useRef } from 'react'

import Box from '@material-ui/core/Box'
import Button from "@material-ui/core/Button";
import Typography from "@material-ui/core/Typography";
import useScrollTrigger from "@material-ui/core/useScrollTrigger";
import SvgIcon from '@material-ui/core/SvgIcon'

//header
import { AppBar, IconButton, Toolbar } from "@material-ui/core";


import Slide from "@material-ui/core/Slide";
import { makeStyles } from "@material-ui/core/styles";

export default function Header (props) {
    const trigger = useScrollTrigger({ threshold: 20 });

    const buttons = ["DevOps", "Backend", "Frontend"];
  
    const color = useToggleAfterInitialRender(props.primary, 'primary', 'textPrimary') 
    const bgColor = useToggleAfterInitialRender(props.primary, 'transparent', 'default')

    const classes = makeStyles((theme) => ({
      mainMenu: {
          justifyContent: "center",
      },
      name: {
        [theme.breakpoints.up("sm")]: { display: 'none'},
      }
    }))();
  
    return (
        <Slide appear={false} in={!trigger} timeout={400}>
          <AppBar color={bgColor} elevation={0}>
            <Toolbar className={classes.mainMenu} >
              <IconButton edge="start" aria-label="menu">
              </IconButton>
              <Box display={{'xs': 'none', 'sm': 'block'}}>
                <Button onClick={props.onClick("About Me")}>
                  <Typography style={{ textTransform: "none" }} color={color} variant="h5">
                    Cooper Melgreen
                  </Typography>
                </Button>
              </Box>
              <Box display={{'xs': 'none', 'sm': 'block'}}>
                <Typography style={{ marginLeft: 20, marginRight: 20 }} color={color} variant="h5">
                  |
                </Typography>
              </Box>
              {buttons.map((text, i) => (
                <Button key={i} onClick={props.onClick(text)}>
                  <Typography className="lower-case" style={{ marginRight: 10 }} color={color} variant="h5">
                    {text}
                  </Typography>
                </Button>
              ))}
            </Toolbar>
          </AppBar>
        </Slide>
    );
}

// TO-DO: check if useReducer is better pattern
const useToggleAfterInitialRender = (bool, value1, value2) => {
  const [state, setState] = useState()
  const firstMount = !useDidMount()

  useEffect(() => {
    console.log(firstMount)
    setState(bool || firstMount ? value1 : value2)
  }, [bool])

  return state
}

function useDidMount() {
  const mountRef = useRef(false);

  useEffect(() => { mountRef.current = true }, []);

  return mountRef.current;
}