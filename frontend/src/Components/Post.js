import React from 'react'
import { useParams } from 'react-router-dom'

import parse from 'html-react-parser'

import { usePostBySlug } from '../API/API'

export default function Post(props) {
    const slug = useParams().slug
    const post = usePostBySlug(slug)

    console.log(post)

    return post.content ? parse(post.content) : null
}