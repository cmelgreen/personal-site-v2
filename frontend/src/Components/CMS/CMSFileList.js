import React from "react";
import { useEffect, useState, useRef } from 'react'
import { Link } from "react-router-dom";

import axios from 'axios'

import { List, ListItem, ListItemText} from '@material-ui/core'

export default function CMSFileList(props) {
  const [selectedIndex, setSelectedIndex] = React.useState();

  const setSelected = i => setSelectedIndex(i)

  const [posts, setPosts] = useState([])

  useEffect(() => {
    axios.get(apiPostSummaries, {params: {numPosts: 10}})
      .then(resp => {
        if ( resp.data.posts ) setPosts(resp.data.posts)
      })
      .catch(() => setPosts([]))
  }, [props.render])
  

  return (
    <List component="nav">
        <ListItem xs={12} sm={12} lg={6} 
          onClick={() => setSelected()}
          button
          component={Link}
          to={"/cms/"}>
          <ListItemText 
            primary="Create New Post"
            primaryTypographyProps={{variant: 'h4'}}
          />
        </ListItem>
        {posts.map((post, i) => (
          <ListItem xs={12} sm={12} lg={6} 
          key={i}
          selected={selectedIndex === i}
          onClick={() => setSelected(i)}
          button
          component={Link}
          to={"/cms/" + post.slug}>
            <ListItemText 
            primary={post.title} 
           primaryTypographyProps={{variant: "h6"}}
            secondary={post.summary}/>
          </ListItem>
        ))}
    </List>
  )
}

const apiPostSummaries = "http://localhost:8080/api/post-summaries"

export const usePostSummaries = (numPosts=10) => {
  const [posts, setPosts] = useState([])

  useEffect(() => {
    axios.get(apiPostSummaries, {params: {numPosts}})
      .then(resp => {
        if ( resp.data.posts ) setPosts(resp.data.posts)
      })
      .catch(() => setPosts([]))
    }, [])

  return posts
}