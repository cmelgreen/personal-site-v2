import React from 'react'

import { useState, useRef, forwardRef } from "react";
import { useInView } from "react-hook-inview";

import Header from './Header'
import Splash from './Splash'
import ContentList from './ContentList'
import AboutMe from './AboutMe'

export default function MainPage(props) {
    const postRef = useRef("postRef");
    const aboutRef = useRef("aboutRef");

    const [category, setCategory] = useState("Posts");
    const [inViewRef, inView] = useInView({ threshold: 0, initialInView: true});
  
    const scrollToPosts = newCategory => () => {
      setCategory(newCategory);
      postRef.current.scrollIntoView({
        behavior: "smooth"
      });
    };

    const scrollToAbout = () => () => {
      setCategory('Posts')
      aboutRef.current.scrollIntoView({
        behavior: "smooth"
      })
    }

    return (
        <>
            <Header className="top-bar" scrollToAbout={scrollToAbout} scrollToPost={scrollToPosts} primary={inView} />
            <Splash ref={inViewRef} timeout={4000}/>
            <AboutMe ref={aboutRef}/>
            <ContentList category={category} ref={postRef}/>
        </>
    )
}