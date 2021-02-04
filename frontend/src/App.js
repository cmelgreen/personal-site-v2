import React from "react";
import { useState, useRef, forwardRef } from "react";

import { BrowserRouter as Router, Route, useHistory  } from 'react-router-dom';

import CssBaseline from "@material-ui/core/CssBaseline";
import { createMuiTheme, ThemeProvider} from '@material-ui/core/styles'

import Box from '@material-ui/core/Box'
import { useInView } from "react-hook-inview";
import "./styles.css";

import firebase from "firebase/app";
import "firebase/auth";
import {
  FirebaseAuthProvider,
  FirebaseAuthConsumer,
  IfFirebaseAuthed,
  IfFirebaseUnAuthed
} from "@react-firebase/auth";
 

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

// Configure Firebase.
const config = require('./credentials/firebase.json')

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
            <FirebaseAuthProvider {...config} firebase={firebase}>
              <Route path={'/cms/:slug?'}>
                <IfFirebaseAuthed>
                  {({ user }) => (
                      <CMS user={user} />
                  )}
                </IfFirebaseAuthed>
                <IfFirebaseUnAuthed>
                  <SignIn />
                </IfFirebaseUnAuthed>
              </Route>
            </FirebaseAuthProvider>
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

const SignIn = (props) => {
  console.log('sign in page')
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const history = useHistory()

  const onSubmit = () => firebase.auth().signInWithEmailAndPassword(email, password)
  .then((userCredential) => {
    firebase.auth().currentUser.getIdToken(/* forceRefresh */ true).then(idToken => {
      console.log(idToken)
    }).catch(error => {
      console.log(error)
    });
    history.push('/cms')
  })
  .catch((error) => {
    console.log(error)
    history.push('/cms')
  });

  return (
  // <form onSubmit={onSubmit}>
  <>
    <label>
      Email:
      <input type="text" value={email} onChange={e => setEmail(e.target.value)} />
    </label>
    <label>
      Password:
      <input type="text" value={password} onChange={e => setPassword(e.target.value)} />
    </label>
    <input type="submit" value="Submit" onClick={onSubmit}/>
  </>
  )
}
