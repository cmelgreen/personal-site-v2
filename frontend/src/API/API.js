import React from 'react'

import { useState, useEffect, useRef } from 'react'
import axios from 'axios';

const apiRoot = "https://api.cmelgreen.com/api/"
export const apiPostSummaries = apiRoot + "post-summaries"
export const apiPost = apiRoot + "post/"

export const usePostSummaries = (numPosts=10) => {
  const [posts, setPosts] = useState([])

  useEffect(() => {
    axios.get(apiPostSummaries, {params: {numPosts}})
      .then(resp => {
        console.log(resp)
        if ( resp.data.posts ) setPosts(resp.data.posts)})
      .catch(() => setPosts([]))
  }, [])

  return posts
}

const api = axios.create({
  baseURL: apiPost,
  validateStatus: (status) => {
      return status == 200;
  }
});

export const createPost = (post, authToken) => {
  return api.post(apiPost, post, {
    headers: {
      'Authorization': `Bearer ${authToken}` 
  }})
}

export const updatePost = (post, authToken) => {
  return api.put(apiPost, post, {
    headers: {
      'Authorization': `Bearer ${authToken}` 
  }})
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

  return post
}

export const deletePost = (post, idToken) => {
  return api.delete(apiPost + post.slug, {
    headers: {
      'Authorization': `Bearer ${idToken}` 
  }})
}