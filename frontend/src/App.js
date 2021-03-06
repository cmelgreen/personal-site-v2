import React, { Suspense } from "react";
import { lazy } from "react";

import { BrowserRouter as Router, Route } from 'react-router-dom';

import CssBaseline from "@material-ui/core/CssBaseline";
import { createMuiTheme, ThemeProvider } from '@material-ui/core/styles'

import "./styles.css";

//content
const CMSWithLogin = lazy(() => import('./Components/CMSWithLogin'))
const MainPage = lazy(() => import('./Components/MainPage'))
const Post = lazy(() => import('./Components/Post'))

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
    fontFamily: 'Frank Ruhl Libre',
    fontSmoothing: 'antialiased',
  },
});

// Configure Firebase.
const config = require('./firebase/firebase.json')

export default function App() {
  // Use react transitions
  return (
    <CssBaseline>
      <ThemeProvider theme={theme}>
        <Suspense fallback={<div>Loading...</div>}>
          <Router>
            <Route path="/" exact component={MainPage} />
            <Route path='/post/:slug' component={Post} />
            <Route path='/cms/:slug?' render={() => <CMSWithLogin config={config}/>}/>
          </Router>
        </Suspense>
      </ThemeProvider>
    </CssBaseline>
  );
}





