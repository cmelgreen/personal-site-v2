import React from 'react'

import { useState, useRef, forwardRef } from "react";
import { useInView } from "react-hook-inview";

import Box from '@material-ui/core/Box'

import Header from './Header'
import Splash from './Splash'
import ContentList from './ContentList'
import AboutMe from './AboutMe'

export default function MainPage(props) {
    const scrollRef = useRef("scrollRef");
    const [category, setCategory] = useState("Content");
    const [inViewRef, inView] = useInView({ threshold: .1, initialInView: true});
  
    const onClick = (newCategory) => () => {
      setCategory(newCategory);
      scrollRef.current.scrollIntoView({
        behavior: "smooth"
      });
    };
    return (
        <>
            <Header className="top-bar" onClick={onClick} primary={inView} />
            <Splash ref={inViewRef} timeout={4000}/>
            <AboutMe />
            <HeaderPadding height={64} ref={scrollRef} /> 
            <ContentList category={category} />
        </>
    )
}

const HeaderPadding = forwardRef((props, ref) => {
    return <Box {...props} ref={ref} />
  })