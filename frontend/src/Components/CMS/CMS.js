import React from "react";
import { useState, useEffect } from 'react'

import axios from 'axios'

import Box from '@material-ui/core/Box'
import Paper from '@material-ui/core/Paper'
import Grid from '@material-ui/core/Grid'
import Container from '@material-ui/core/Container'

import { makeStyles } from "@material-ui/core/styles";

import Editor from './Editor'
import CMSCardList from './CMSFileList'

export default function CMS(props) {
  const classes = useStyles()

  return (
    <Container className={classes.fullscreen}>
        <Grid className={classes.fullscreen} container directon="row" >
          <Grid className={classes.widget} component={Paper} elevation={2} square={true} item xs={3}>
            <Box ml="5em" mt="5em" >
              <CMSCardList posts={props.posts}/>
            </Box>
          </Grid>
          <Grid className="{classes.widget}" item xs={9}>
            <Box component={Paper} elevation={4} square={true} mx="5em" mt="2em">
              <Editor />
            </Box>  
          </Grid>
        </Grid>
    </Container>
  )
}

const useStyles = makeStyles((theme) => ({
  fullscreen: {
    margin: 0,
    padding: 0,
    height: '100vh',
    width: '100vw',
  },
  widget: {
    background: 'linear-gradient(126deg, rgba(255,226,242,1) 7%, rgba(245,210,237,1) 99%)'
  }
}))

