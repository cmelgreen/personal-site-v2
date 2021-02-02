import React from 'react'

import { useState, useEffect, useRef } from 'react'
import { TextField } from '@material-ui/core'
import MUIRichTextEditor from 'mui-rte';
import { useParams, useHistory } from 'react-router-dom'
import { convertToRaw } from 'draft-js'

import { makeStyles } from "@material-ui/core/styles";

// import { usePostByID, usePostSummaries, createPost, updatePost, deletePost } from '../../../Utils/ContentAPI'

import Button from '@material-ui/core/Button';
import { List, ListItem, ListItemText} from '@material-ui/core'
import Grid from '@material-ui/core/Grid';

import ClearIcon from '@material-ui/icons/Clear';

import axios from 'axios'

export default function Editor(props) {
  const newPost = {title: '', slug: '', content: '', tags: []}
  const history = useHistory()

  const classes = useStyles() 

  const slug = useParams().slug
  const [post, setPost] = useState(newPost)

  const firstLoad = useRef(false)

  useEffect(() => {
    if ( slug ) {
      axios.get(apiPost + slug, {params: {raw: true}})
        .then(resp => {setPost(resp.data)})
        .catch(resp => {setPost(newPost)})
    }

  }, [slug])

  const save = useRef(false)

  useEffect(() => {
    if ( save.current ) {
      if ( slug ) {
        updatePost(post)
      } else {
        createPost(post)
          .then(() => history.push('/cms/' + post.slug))
          .catch(e => console.log(e))
      }

      save.current = false
    }
  }, [save.current, post])


  const postFields = ['title', 'slug', 'summary']

  const onSave = (richText) => {
    setPost({...post, content: richText})
    save.current = true
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
        <Grid item key={i}>
          <TextField
            id={field}
            label={field}
            InputLabelProps={{ shrink: true }} 
            value={post[field] ? post[field] : ""}
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
        //   () => {
        //   console.log('Content', post.content)
        //   return post.content ? post.content : ''
        // }}
        //onChange={handleChange}
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

const api = axios.create({
  baseURL: apiPost,
  validateStatus: (status) => {
       return status == 200;
   },
});

export const createPost = post => {
  return api.post(apiPost, post)
}

export const updatePost = post => {
  return api.put(apiPost, post)
}

export const usePostBySlug = (slug, raw=false) => {
  const [post, setPost] = useState({})

  useEffect(() => {
    if (slug) {
      api.get(apiPost + slug, {params: {raw}})
        .then(resp => setPost(resp.data))
        .catch(resp => setPost({}))
    }
  }, [slug])

  return [post, setPost]
}

export const deletePost = (post) => {
  return api.delete(apiPost + post.slug)
}
