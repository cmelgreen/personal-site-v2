import React, { Suspense } from "react";
import { lazy } from "react";

import { BrowserRouter as Router, Route } from 'react-router-dom';

import CssBaseline from "@material-ui/core/CssBaseline";
import { createMuiTheme, ThemeProvider } from '@material-ui/core/styles'

import "./styles.css";

//content
const CMSWithLogin = lazy(() => import('./Components/CMSWithLogin'))
const MainPage = lazy(() => import('./Components/MainPage'))

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

// Configure Firebase.
const config = require('./firebase/firebase.json')

export default function App() {
  // Use react transitions
  return (
    <div className="App">
      <CssBaseline>
        <ThemeProvider theme={theme}>
          <Suspense fallback={Loading}>
            <Router>
              <Route path={'/cms/:slug?'}>
                <CMSWithLogin config={config}/>
              </Route>
              <Route path="/" exact={true}>
                <MainPage />
              </Route>
            </Router>
          </Suspense>
        </ThemeProvider>
      </CssBaseline>
    </div>
  );
}

const Loading = () => (<div>Loading...</div>)




