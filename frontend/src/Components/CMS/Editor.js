import React from 'react'

import { useState, useEffect } from 'react'
import { TextField } from '@material-ui/core'
import MUIRichTextEditor from 'mui-rte';
import { useParams, useHistory } from 'react-router-dom'

import { makeStyles } from "@material-ui/core/styles";

// import { usePostByID, usePostSummaries, createPost, updatePost, deletePost } from '../../../Utils/ContentAPI'

import Button from '@material-ui/core/Button';
import { List, ListItem, ListItemText} from '@material-ui/core'
import Grid from '@material-ui/core/Grid';

import ClearIcon from '@material-ui/icons/Clear';

import axios from 'axios'

export default function Editor(props) {
  const history = useHistory()

  const classes = useStyles() 

  const slug = useParams().slug
  const [post, setPost] = useState({})

  useEffect(() => {
    if (slug) {
      axios.get(apiPost + slug, {params: {raw: true}})
        .then(resp => setPost(resp.data))
        .catch(resp => setPost({}))
    }
  }, [slug])

  const postFields = ['title', 'slug', 'summary']

  console.log('Post:', post)
  
  const onSave = (content) => {
    setPost({...post, content: content})
    slug ? updatePost(post) : createPost(post)
    console.log(post)
    history.push("/cms/" + post.slug) 
  }

  // const post = usePostByID(useParams().postID, true)

  // const [id, setID] = useState(post.id)
  // const [title, setTitle] = useState(post.title)
  // const [summary, setSummary] = useState(post.summary)
  // const [tags, setTags] = useState(post.tags)

  // useEffect(() => {
  //   setID(post.id)
  //   setTitle(post.title)
  //   setSummary(post.summary)
  // }, [post])

  // const [saveState, setSaveState] = useState(true)
  // usePostSummaries(-1, saveState)


  // const onSave = (data) => {
  //   saveType(data)
  //   setSaveState(!saveState)
  //   history.push("/cms/"+title) 
  // }

  return (
    <Grid container className={classes.fullscreen} direction="column" spacing={2}>
      {postFields.map((field, i) => (
        <Grid item>
          <TextField
            id={field}
            label={field}
            InputLabelProps={{ shrink: true }} 
            value={post[field]}
            onChange={e => setPost({...post, [field]: e.target.value})}
            variant="outlined"
            fullWidth={true}
          />
        </Grid>
      ))}
        <Grid item>
          <TextField
            label="tags"
            InputLabelProps={{ shrink: true }} 
            value={post.tags ? post.tags : ""}
            onChange={e => setPost({...post, tags: e.target.value.split(',')})}
            variant="outlined"
            fullWidth={true}
          />
        </Grid>
      <Grid item>
        <Button
          variant="outlined"
          component="label"
        >
          Upload File
          <input
            type="file"
            hidden
          />
      </Button>
      </Grid>
      <Grid item>
        <MUIRichTextEditor 
        defaultValue={post.content} 
        onSave={onSave} 
        controls={["title", "bold", "italic", "underline", "strikethrough", "highlight", "undo", "redo", "link", "media", "numberList", "bulletList", "quote", "code", "clear", "save", "deletePost"]}
        customControls={[
          {
              name: "deletePost",
              icon: <ClearIcon />,
              type: "callback",
              onClick: () => {deletePost(post); history.push("/cms/")}
          }
        ]}
        />
      </Grid>
    </Grid>
  );
}

const useStyles = makeStyles((theme) => ({
  fullscreen: {
    margin: 0,
    padding: 0,
    height: '100vh',
    width: '100%',
  },
}))

const apiPost = 'http://localhost:8080/api/post/'

export const createPost = post => {
  axios.post(apiPost, post)
    .then(resp => console.log('Created', resp))
    .catch(resp => console.log('Error creating post', resp))
}

export const updatePost = post => {
  axios.put(apiPost, post)
    .then(resp => console.log('Updated', resp))
    .catch(resp => console.log('Error updating post', resp))
}

export const usePostBySlug = (slug, raw=false) => {
  const [post, setPost] = useState({})

  useEffect(() => {
    if (slug) {
      axios.get(apiPost + slug, {params: {raw}})
        .then(resp => setPost(resp.data))
        .catch(resp => setPost({}))
    }
  }, [slug])

  console.log('Postbyslug:',  post)

  return [post, setPost]
}

export const deletePost = (post) => {
  axios.delete(apiPost + post.slug)
    .then(resp => console.log('Deleted'))
    .catch(resp => console.log('Error deleting post', resp))
}
