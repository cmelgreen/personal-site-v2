import React from "react";
import { useState, useRef, forwardRef } from "react";

import CssBaseline from "@material-ui/core/CssBaseline";
import Box from '@material-ui/core/Box'

import { useInView } from "react-hook-inview";
import "./styles.css";

//import Grow from '@material-ui/core/Grow';


//content
import { Splash } from "./Splash.js";
import ContentList from "./ContentList.js";
import { Header } from './Header.js'

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
    <div className="App" >
      <CssBaseline>
        <Header className="top-bar" onClick={onClick} primary={inView} />
        <Splash ref={inViewRef} />
        <HeaderPadding height={64} ref={scrollRef} /> 
        <ContentList category={category} />
      </CssBaseline>
    </div>
  );
}

const HeaderPadding = forwardRef((props, ref) => {
  return <Box {...props} ref={ref} />
})