import React from "react";
import { useEffect, useState, useRef } from 'react'
import { Link } from "react-router-dom";

import axios from 'axios'

import { usePostSummaries, apiPostSummaries } from '../../API/API'

import { Button, List, ListItem, ListItemText } from '@material-ui/core'

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

  const uploadImage = (image) => {
    const formData = new FormData()
    const reader = new FileReader()
    const url = reader.readAsDataURL(image)

    formData.append(
      "image",
      image
    )

    axios.post('https://api.cmelgreen.com/api/img/static/', formData, {
      headers: { 'content-type': 'multipart/form-data' }
    })
    .then(resp => {
      console.log(resp)
    })
    .catch(resp => console.log(resp))
  }
  

  return (
    <List component="nav">
        <ListItem xs={12} sm={12} lg={6} 
          button
          component="label"
        >
            <input type="file" onChange={e => uploadImage(e.target.files[0])} hidden></input>
          <ListItemText 
            primary="Upload Static File"
            primaryTypographyProps={{variant: 'h4'}}
          />
        </ListItem>
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