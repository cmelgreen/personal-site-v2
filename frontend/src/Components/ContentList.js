import React, { useState, useEffect } from 'react';
import { forwardRef } from 'react';

import Container from '@material-ui/core/Container'
import Typography from '@material-ui/core/Typography';
import ContentCard from './ContentCard.js'
import Grid from '@material-ui/core/Grid'

import { makeStyles } from "@material-ui/core/styles";

import { usePostSummaries } from '../API/API'

const ContentList = forwardRef((props, ref) => {
  const posts = usePostSummaries()

  const useStyles = makeStyles((theme) => ({
    categoryTitle: {
      margin: '5% auto 5% auto',
    },
    contentGrid: {
      minHeight: '100vh',
      justifyContent: 'center',
      [theme.breakpoints.up("xs")]: {margin: '0% 5% 0% 5%'},
      [theme.breakpoints.up("sm")]: {margin: '0% 10% 0% 10%'},
      [theme.breakpoints.up("lg")]: {margin: '0% 15% 0% 15%'},
      [theme.breakpoints.up('xl')]: {margin: '0% 25% 0% 25%'},
      // marginLeft: '10%',
      // marginRight: '10%',
    }
  }))

  const classes = useStyles()

  const filterContent = post => {
    console.log(post.category, props.category)
    return post.category === props.category
  }

  const filteredPosts = ( props.category === 'Posts') ?
    posts :
    posts.filter(filterContent)

  return (
    <>
      <Container className={classes.categoryTitle} padding={'10%'}>
        <Typography ref={ref} variant='h1' align='center'>{props.category}</Typography>
      </Container>      
      <Grid className={classes.contentGrid} alignContent='center'>
        <Grid container spacing={2} alignItems='center'>
          {filteredPosts.map((post, i) => 
            <Grid item xs={12} sm={12} md={6} lg={6} >
              <ContentCard key={i} post={post}/>
            </Grid>
          )}
        </Grid>
      </Grid>
    </>
  )
})

export default ContentList