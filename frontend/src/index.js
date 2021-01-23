import React from "react";
import ReactDOM from "react-dom";

import App from "./App";

import { createMuiTheme, ThemeProvider } from '@material-ui/core/styles'

const rootElement = document.getElementById("root");

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
    h6: {
      fontStyle: 'italic',
    },
    button: {
      fontWeight: 600,
      color: 'secondary',
    },
  },
});

ReactDOM.render(
  <React.StrictMode>
    <ThemeProvider theme={theme}>
      <App />
    </ThemeProvider>
  </React.StrictMode>,
  rootElement
);