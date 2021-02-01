import React from "react";
import { useState, useRef, forwardRef } from "react";

import { BrowserRouter as Router, Route } from 'react-router-dom';

import CssBaseline from "@material-ui/core/CssBaseline";
import { createMuiTheme, ThemeProvider} from '@material-ui/core/styles'

import Box from '@material-ui/core/Box'
import { useInView } from "react-hook-inview";
import "./styles.css";

//content
import { Splash } from "./Components/Splash.js";
import ContentList from "./Components/ContentList.js";
import { Header } from './Components/Header.js'
import CMS from "./Components/CMS/CMS"

const theme = createMuiTheme({
  palette: {
    type: 'light',
    primary: {
      main: '#FFF',
    },
    secondary: {
      main: '#fe5b25',
    },
    textPrimary: {
      main: '#000',
    },
    background: {
      default: 'white'
    }
  },
  typography: {
    //fontFamily: 'Karla',
    //fontFamily: 'Vollkorn',
    //fontFamily: 'Lora',
    fontFamily: 'Frank Ruhl Libre',
    fontSmoothing: 'antialiased',
    // button: {
    //   fontWeight: 600,
    //   color: 'secondary',
    // },
  },
});

export default function App() {
  const scrollRef = useRef("scrollRef");
  const [category, setCategory] = useState("Content");
  const [inViewRef, inView] = useInView({ threshold: .1, initialInView: true});

  const onClick = (newCategory) => () => {
    setCategory(newCategory);
    scrollRef.current.scrollIntoView({
      behavior: "smooth"
    });
  };

  // Use react transitions
  return (
    <div className="App">
      <CssBaseline>
        <ThemeProvider theme={theme}>
          <Router>
          <Route path='/cms' exact={true}>
              <CMS />
            </Route>
            <Route path='/cms/:slug'>
              <CMS />
            </Route>
            <Route path="/" exact={true}>
              <Header className="top-bar" onClick={onClick} primary={inView} />
              <Splash ref={inViewRef} />
              <HeaderPadding height={64} ref={scrollRef} /> 
              <ContentList category={category} />
            </Route>
          </Router>
        </ThemeProvider>
      </CssBaseline>
    </div>
  );
}

const HeaderPadding = forwardRef((props, ref) => {
  return <Box {...props} ref={ref} />
})

